package types

// AccountWithTransactions interface is a wrapper of an account and its transactions
type AccountWithTransactions interface {
	Account
	GetTotalTransactionCount() int
	GetTransactions() []interface{}
}
