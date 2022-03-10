package provider

import (
	"math/big"

	"github.com/pokt-foundation/pocket-go/pkg/models"
)

// Provider interface that represents a provider
type Provider interface {
	// Account
	GetBalance(address string) big.Int
	GetTransactionCount(address string) int
	GetType(address string) AddressType
	// TXs
	SendTransaction(signerAddress string, signedTransaction string) models.TransactionResponse
	// Network
	GetBlock(blockNumber string) models.Block
	GetTransaction(transactionHash string) models.TransactionResponse
	GetBlockNumber() int
	GetNodes(getNodesOptions models.GetNodesOptions) []models.Node
	GetNode(address string, options models.GetNodesOptions) models.Node
	GetApps(getAppOptions models.GetAppOptions) []models.App
	GetApp(address string, options models.GetAppOptions) models.App
	GetAccount(address string) models.Account
	GetAccountWithTransactions(address string) models.AccountWithTransactions
	// TODO: Add methods for params/requestChallenge
}
