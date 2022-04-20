package transactionbuilder

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"math/big"

	"github.com/pokt-foundation/pocket-go/pkg/provider"
	"github.com/pokt-foundation/pocket-go/pkg/signer"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	coreTypes "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
)

const defaultTXFee = int64(10000)

var (
	// ErrNoSigner error when no signer is provided
	ErrNoSigner = errors.New("no signer provided")
	// ErrNoProvider error when no provider is provided
	ErrNoProvider = errors.New("no provider provided")
	// ErrNoChainID error when no chain ID is provided
	ErrNoChainID = errors.New("no chain id provided")
	// ErrNoTxMsg error when no Tx Msg is provided
	ErrNoTxMsg = errors.New("no tx msg provided")
)

// PocketTransactionBuilder represents implementation of transaction builder package
type PocketTransactionBuilder struct {
	provider provider.Provider
	signer   signer.Signer
}

// NewPocketTransactionBuilder returns an instance of PocketTransactionBuilder
func NewPocketTransactionBuilder(provider provider.Provider, signer signer.Signer) *PocketTransactionBuilder {
	return &PocketTransactionBuilder{
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

func (t *PocketTransactionBuilder) validateTransactionRequest(chainID string, txMsg TxMsg) error {
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
		return ErrNoTxMsg
	}

	return nil
}

func (t *PocketTransactionBuilder) signTransaction(chainID, memo, coinDenom string, fee int64, txMsg TxMsg) (string, error) {
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

	signature, err := t.signer.GetKeyManager().SignBytes(signBytes)
	if err != nil {
		return "", err
	}

	publicKey, err := crypto.NewPublicKey(t.signer.GetKeyManager().GetPublicKey())
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
func (t *PocketTransactionBuilder) CreateTransaction(chainID string, txMsg TxMsg, options *TransactionOptions) (*provider.SendTransactionInput, error) {
	err := t.validateTransactionRequest(chainID, txMsg)
	if err != nil {
		return nil, err
	}

	memo, coinDenom, fee := getOptionalParams(options)

	signedTX, err := t.signTransaction(chainID, memo, coinDenom, fee, txMsg)
	if err != nil {
		return nil, err
	}

	return &provider.SendTransactionInput{
		Address:     t.signer.GetKeyManager().GetAddress(),
		RawHexBytes: signedTX,
	}, nil
}

// Submit does the transaction from raw input
func (t *PocketTransactionBuilder) Submit(chainID string, txMsg TxMsg, options *TransactionOptions) (*provider.SendTransactionOutput, error) {
	sendTransactionInput, err := t.CreateTransaction(chainID, txMsg, options)
	if err != nil {
		return nil, err
	}

	return t.provider.SendTransaction(sendTransactionInput)
}
