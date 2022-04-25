package transactionbuilder

import (
	"github.com/pokt-foundation/pocket-go/pkg/provider"
	coreTypes "github.com/pokt-network/pocket-core/types"
)

// TransactionBuilder interface that represents functionalities of the transaction builder package
type TransactionBuilder interface {
	Submit(chainID ChainID, txMsg TxMsg, options *TransactionOptions) (*provider.SendTransactionOutput, error)
	CreateTransaction(chainID ChainID, txMsg TxMsg, options *TransactionOptions) (*provider.SendTransactionInput, error)
}

// TxMsg interface that represents message to be sent as transaction
type TxMsg interface {
	coreTypes.ProtoMsg
}

// Provider interface representing provider functions necessary for Transaction Builder Package
type Provider interface {
	SendTransaction(input *provider.SendTransactionInput) (*provider.SendTransactionOutput, error)
}

// Signer interface representing signer functions necessary for Transaction Builder package
type Signer interface {
	SignBytes(payload []byte) ([]byte, error)
	GetAddress() string
	GetPublicKey() string
}
