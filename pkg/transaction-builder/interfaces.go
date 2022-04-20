package transactionbuilder

import (
	"github.com/pokt-foundation/pocket-go/pkg/provider"
	coreTypes "github.com/pokt-network/pocket-core/types"
)

// TransactionBuilder interface that represents functionalities of the transaction builder package
type TransactionBuilder interface {
	Submit(chainID, memo string, fee int64, txMsg TxMsg, coinDenom CoinDenom) (*provider.SendTransactionOutput, error)
	CreateTransaction(chainID, memo string, fee int64, txMsg TxMsg, coinDenom CoinDenom) (*provider.SendTransactionInput, error)
}

// TxMsg interface that represents message to be sent as transaction
type TxMsg interface {
	coreTypes.ProtoMsg
}
