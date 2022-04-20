package transactionbuilder

import (
	"github.com/pokt-foundation/pocket-go/pkg/provider"
	coreTypes "github.com/pokt-network/pocket-core/types"
)

// TransactionBuilder interface that represents functionalities of the transaction builder package
type TransactionBuilder interface {
	Submit(chainID string, txMsg TxMsg, options *TransactionOptions) (*provider.SendTransactionOutput, error)
	CreateTransaction(chainID string, txMsg TxMsg, options *TransactionOptions) (*provider.SendTransactionInput, error)
}

// TxMsg interface that represents message to be sent as transaction
type TxMsg interface {
	coreTypes.ProtoMsg
}
