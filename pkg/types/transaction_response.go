package types

// TransactionResponse interface that represents response of a transaction
type TransactionResponse interface {
	GetCode() int
	GetCodeSpace() string
	GetData() string
	GetHash() string
	GetHeight() int
	GetInfo() string
	GetRawLog() string
	GetTimestamp() string
	GetTx() string
}
