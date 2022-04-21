package transactionbuilder

// TransactionOptions represents optional parameters for transaction request
type TransactionOptions struct {
	Memo      string
	Fee       int64
	CoinDenom CoinDenom
}
