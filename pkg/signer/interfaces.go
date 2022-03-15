package signer

import "github.com/pokt-foundation/pocket-go/pkg/provider"

// TransactionRequest not implemented
// TODO: implement Transaction Request
type TransactionRequest interface{}

// Signer interface that represents a Signer implementation
type Signer interface {
	Connect(requestProvider provider.Provider)
	SignTransaction(req TransactionRequest) (string, error)
	GetKeyManager() *KeyManager
	GetProvider() provider.Provider
	IsSigner() bool
}
