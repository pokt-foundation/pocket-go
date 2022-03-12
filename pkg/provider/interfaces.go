package provider

import (
	"math/big"
)

// Provider interface that represents a provider
type Provider interface {
	// Account
	GetBalance(address string) (*big.Int, error)
	GetTransactionCount(address string) (int, error)
	GetType(address string) (AddressType, error)
	// TXs
	SendTransaction(signerAddress string, signedTransaction string) (*SendTransactionReponse, error)
	// Network
	GetBlock(blockNumber int) (*GetBlockResponse, error)
	GetTransaction(transactionHash string) (*GetTransactionReponse, error)
	GetBlockNumber() (int, error)
	GetNodes(height int, options *GetNodesOptions) (*GetNodesResponse, error)
	GetNode(address string, options *GetNodeOptions) (*GetNodeResponse, error)
	GetApps(height int, options *GetAppsOptions) (*GetAppsResponse, error)
	GetApp(address string, options *GetAppOptions) (*GetAppResponse, error)
	GetAccount(address string) (*GetAccountResponse, error)
	GetAccountWithTransactions(address string) (*GetAccountWithTransactionsResponse, error)
	// TODO: Add methods for params/requestChallenge
}
