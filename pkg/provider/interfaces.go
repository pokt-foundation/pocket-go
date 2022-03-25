package provider

import (
	"math/big"
)

// Provider interface that represents a provider
type Provider interface {
	// Account
	GetBalance(address string) (*big.Int, error)
	GetTransactionCount(address string) (int, error)
	GetType(address string, options *GetTypeOptions) (AddressType, error)
	// TXs
	SendTransaction(signerAddress string, signedTransaction string) (*SendTransactionResponse, error)
	// Network
	GetBlock(blockNumber int) (*GetBlockResponse, error)
	GetTransaction(transactionHash string) (*GetTransactionResponse, error)
	GetBlockNumber() (int, error)
	GetNodes(height int, options *GetNodesOptions) (*GetNodesResponse, error)
	GetNode(address string, options *GetNodeOptions) (*GetNodeResponse, error)
	GetApps(height int, options *GetAppsOptions) (*GetAppsResponse, error)
	GetApp(address string, options *GetAppOptions) (*GetAppResponse, error)
	GetAccount(address string) (*GetAccountResponse, error)
	GetAccountWithTransactions(address string) (*GetAccountWithTransactionsResponse, error)
	Dispatch(appPublicKey, chain string, sessionHeight int, options *DispatchRequestOptions) (*DispatchResponse, error)
	Relay(rpcURL string, input *Relay, options *RelayRequestOptions) (string, error)
	// TODO: Add methods for params/requestChallenge
}
