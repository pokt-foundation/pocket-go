package transactionbuilder

// CoinDenom enum that represents all coin denominations of Pocket
type CoinDenom string

const (
	// Upokt represents upokt denomination
	Upokt CoinDenom = "upokt"
	// Pokt represents pokt denomination
	Pokt CoinDenom = "pokt"
)

// ChainID enum that represents possible chain IDs for transactions
type ChainID string

const (
	// Mainnet use for connecting to Pocket Mainnet
	Mainnet ChainID = "mainnet"
	// Testnet use for connecting to Pocket Testnet
	Testnet ChainID = "testnet"
	// Localnet use for connecting to Pocket Localnet
	Localnet ChainID = "localnet"
)
