package transactionbuilder

import (
	"crypto/rand"
	"encoding/hex"
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

// CreateTransaction returns input necessary for doing a transaction
func (t *PocketTransactionBuilder) CreateTransaction(chainID, memo string, fee int64, txMsg TxMsg, coinDenom CoinDenom) (*provider.SendTransactionInput, error) {
	feeStruct := coreTypes.Coins{
		coreTypes.Coin{
			Amount: coreTypes.NewInt(fee),
			Denom:  string(coinDenom),
		},
	}

	entropy, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	signBytes, err := auth.StdSignBytes(chainID, entropy.Int64(), feeStruct, txMsg, memo)
	if err != nil {
		return nil, err
	}

	signature, err := t.signer.GetKeyManager().SignBytes(signBytes)
	if err != nil {
		return nil, err
	}

	publicKey, err := crypto.NewPublicKey(t.signer.GetKeyManager().GetPublicKey())
	if err != nil {
		return nil, err
	}

	signatureStruct := authTypes.StdSignature{PublicKey: publicKey, Signature: signature}

	tx := authTypes.NewTx(txMsg, feeStruct, signatureStruct, memo, entropy.Int64())

	txBytes, err := auth.DefaultTxEncoder(app.Codec())(tx, -1)
	if err != nil {
		return nil, err
	}

	return &provider.SendTransactionInput{
		Address:     t.signer.GetKeyManager().GetAddress(),
		RawHexBytes: hex.EncodeToString(txBytes),
	}, nil
}

// Submit does the transaction from raw input
func (t *PocketTransactionBuilder) Submit(chainID, memo string, fee int64, txMsg TxMsg, coinDenom CoinDenom) (*provider.SendTransactionOutput, error) {
	sendTransactionInput, err := t.CreateTransaction(chainID, memo, fee, txMsg, coinDenom)
	if err != nil {
		return nil, err
	}

	return t.provider.SendTransaction(sendTransactionInput)
}
