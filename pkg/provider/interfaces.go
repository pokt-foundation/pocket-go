package provider

import (
	"math/big"
)

// Provider interface that represents a provider
type Provider interface {
	// Account
	GetBalance(address string, options *GetBalanceOptions) (*big.Int, error)
	GetTransactionCount(address string, options *GetTransactionCountOptions) (int, error)
	GetType(address string, options *GetTypeOptions) (AddressType, error)
	// TXs
	SendTransaction(input *SendTransactionInput) (*SendTransactionOutput, error)
	// Network
	GetBlock(blockNumber int) (*GetBlockOutput, error)
	GetTransaction(transactionHash string, options *GetTransactionOptions) (*GetTransactionOutput, error)
	GetBlockHeight() (int, error)
	GetNodes(height int, options *GetNodesOptions) (*GetNodesOutput, error)
	GetNode(address string, options *GetNodeOptions) (*GetNodeOutput, error)
	GetApps(height int, options *GetAppsOptions) (*GetAppsOutput, error)
	GetApp(address string, options *GetAppOptions) (*GetAppOutput, error)
	GetAccount(address string, options *GetAccountOptions) (*GetAccountOutput, error)
	GetAccountTransactions(address string, options *GetAccountTransactionsOptions) (*GetAccountTransactionsOutput, error)
	Dispatch(appPublicKey, chain string, options *DispatchRequestOptions) (*DispatchOutput, error)
	Relay(rpcURL string, input *Relay, options *RelayRequestOptions) (*RelayOutput, error)
	// TODO: Add methods for params/requestChallenge
}
