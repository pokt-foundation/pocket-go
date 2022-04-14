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
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
)

type PocketTransactionBuilder struct {
	provider provider.Provider
	signer   signer.Signer
}

func NewPocketTransactionBuilder(provider provider.Provider, signer signer.Signer) *PocketTransactionBuilder {
	return &PocketTransactionBuilder{
		provider: provider,
		signer:   signer,
	}
}

func (t *PocketTransactionBuilder) CreateTransaction(chainID, memo string, fee int64, txMsg TxMsg, coinDenom provider.CoinDenom) (*provider.SendTransactionInput, error) {
	feeStruct := types.Coins{
		types.Coin{
			Amount: types.NewInt(fee),
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

	signature, err := t.signer.GetKeyManager().Sign(signBytes)
	if err != nil {
		return nil, err
	}

	publicKey, err := crypto.NewPublicKey(t.signer.GetKeyManager().GetPublicKey())
	if err != nil {
		return nil, err
	}

	signatureStruct := authTypes.StdSignature{PublicKey: publicKey, Signature: []byte(signature)}

	tx := authTypes.NewTx(txMsg, feeStruct, signatureStruct, memo, entropy.Int64())

	txBytes, err := auth.DefaultTxEncoder(app.Codec())(tx, -1)
	if err != nil {
		return nil, err
	}

	return &provider.SendTransactionInput{
		Adress:      t.signer.GetKeyManager().GetAddress(),
		RawHexBytes: hex.EncodeToString(txBytes),
	}, nil
}

func (t *PocketTransactionBuilder) Submit(chainID, memo string, fee int64, txMsg TxMsg, coinDenom provider.CoinDenom) (*provider.SendTransactionOutput, error) {
	sendTransactionInput, err := t.CreateTransaction(chainID, memo, fee, txMsg, coinDenom)
	if err != nil {
		return nil, err
	}

	return t.provider.SendTransaction(sendTransactionInput)
}
