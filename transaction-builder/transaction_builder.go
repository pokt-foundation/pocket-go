// Package transactionbuilder is a helper for doing transactions with simpler input that just using the package Provider
// Underneath uses the package Provider
package transactionbuilder

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"math/big"

	"github.com/pokt-foundation/pocket-go/provider"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	coreTypes "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
)

const defaultTXFee = int64(10000)

// CoinDenom enum that represents all coin denominations of Pocket
type CoinDenom string

const (
	// Upokt represents upokt denomination
	Upokt CoinDenom = "upokt"
	// Pokt represents pokt denomination
	Pokt CoinDenom = "pokt"
)

// ChainID enum that represents possible chain IDs for transactions
type ChainID string

const (
	// Mainnet use for connecting to Pocket Mainnet
	Mainnet ChainID = "mainnet"
	// Testnet use for connecting to Pocket Testnet
	Testnet ChainID = "testnet"
	// Localnet use for connecting to Pocket Localnet
	Localnet ChainID = "localnet"
)

var (
	// ErrNoSigner error when no signer is provided
	ErrNoSigner = errors.New("no signer provided")
	// ErrNoProvider error when no provider is provided
	ErrNoProvider = errors.New("no provider provided")
	// ErrNoChainID error when no chain ID is provided
	ErrNoChainID = errors.New("no chain id provided")
	// ErrNoTransactionMessage error when no Transaction Message is provided
	ErrNoTransactionMessage = errors.New("no transaction message provided")
)

// Provider interface representing provider functions necessary for Transaction Builder Package
type Provider interface {
	SendTransactionWithCtx(ctx context.Context, input *provider.SendTransactionInput) (*provider.SendTransactionOutput, error)
}

// Signer interface representing signer functions necessary for Transaction Builder package
type Signer interface {
	SignBytes(payload []byte) ([]byte, error)
	GetAddress() string
	GetPublicKey() string
}

// TransactionBuilder represents implementation of transaction builder package
type TransactionBuilder struct {
	provider Provider
	signer   Signer
}

// TransactionOptions represents optional parameters for transaction request
type TransactionOptions struct {
	Memo      string
	Fee       int64
	CoinDenom CoinDenom
}

// NewTransactionBuilder returns an instance of TransactionBuilder
func NewTransactionBuilder(provider Provider, signer Signer) *TransactionBuilder {
	return &TransactionBuilder{
		provider: provider,
		signer:   signer,
	}
}

func getOptionalParams(options *TransactionOptions) (string, string, int64) {
	memo := ""
	coinDenom := Upokt
	fee := defaultTXFee

	if options != nil {
		memo = options.Memo

		if options.CoinDenom != "" {
			coinDenom = options.CoinDenom
		}

		if options.Fee != 0 {
			fee = options.Fee
		}
	}

	return memo, string(coinDenom), fee
}

func (t *TransactionBuilder) validateTransactionRequest(chainID ChainID, txMsg TransactionMessage) error {
	if t.provider == nil {
		return ErrNoProvider
	}

	if t.signer == nil {
		return ErrNoSigner
	}

	if chainID == "" {
		return ErrNoChainID
	}

	if txMsg == nil {
		return ErrNoTransactionMessage
	}

	return nil
}

func (t *TransactionBuilder) signTransaction(chainID, memo, coinDenom string, fee int64, txMsg TransactionMessage) (string, error) {
	feeStruct := coreTypes.Coins{
		coreTypes.Coin{
			Amount: coreTypes.NewInt(fee),
			Denom:  coinDenom,
		},
	}

	entropy, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "", err
	}

	signBytes, err := auth.StdSignBytes(chainID, entropy.Int64(), feeStruct, txMsg, memo)
	if err != nil {
		return "", err
	}

	signature, err := t.signer.SignBytes(signBytes)
	if err != nil {
		return "", err
	}

	publicKey, err := crypto.NewPublicKey(t.signer.GetPublicKey())
	if err != nil {
		return "", err
	}

	signatureStruct := authTypes.StdSignature{PublicKey: publicKey, Signature: signature}

	tx := authTypes.NewTx(txMsg, feeStruct, signatureStruct, memo, entropy.Int64())

	txBytes, err := auth.DefaultTxEncoder(app.Codec())(tx, -1)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(txBytes), nil
}

// CreateTransaction returns input necessary for doing a transaction
func (t *TransactionBuilder) CreateTransaction(chainID ChainID, txMsg TransactionMessage, options *TransactionOptions) (*provider.SendTransactionInput, error) {
	err := t.validateTransactionRequest(chainID, txMsg)
	if err != nil {
		return nil, err
	}

	memo, coinDenom, fee := getOptionalParams(options)

	signedTX, err := t.signTransaction(string(chainID), memo, coinDenom, fee, txMsg)
	if err != nil {
		return nil, err
	}

	return &provider.SendTransactionInput{
		Address:     t.signer.GetAddress(),
		RawHexBytes: signedTX,
	}, nil
}

// Submit does the transaction from raw input
func (t *TransactionBuilder) Submit(chainID ChainID, txMsg TransactionMessage, options *TransactionOptions) (*provider.SendTransactionOutput, error) {
	return t.SubmitWithCtx(context.Background(), chainID, txMsg, options)
}

// SubmitWithCtx does the transaction from raw input
func (t *TransactionBuilder) SubmitWithCtx(ctx context.Context, chainID ChainID, txMsg TransactionMessage, options *TransactionOptions) (*provider.SendTransactionOutput, error) {
	sendTransactionInput, err := t.CreateTransaction(chainID, txMsg, options)
	if err != nil {
		return nil, err
	}

	return t.provider.SendTransactionWithCtx(ctx, sendTransactionInput)
}
