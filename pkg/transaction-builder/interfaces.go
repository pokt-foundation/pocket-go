package transactionbuilder

import (
	"github.com/pokt-foundation/pocket-go/pkg/provider"
	"github.com/pokt-network/pocket-core/types"
)

type TransactionBuilder interface {
	Submit(chainID, memo string, fee int64, txMsg TxMsg, coinDenom provider.CoinDenom) (*provider.SendTransactionOutput, error)
	CreateTransaction(chainID, memo string, fee int64, txMsg TxMsg, coinDenom provider.CoinDenom) (*provider.SendTransactionInput, error)
}

type TxMsg interface {
	types.ProtoMsg
}
