package signer

import (
	"errors"

	"github.com/pokt-foundation/pocket-go/pkg/provider"
)

var (
	// ErrNotImplemented error when function is not implemented
	ErrNotImplemented = errors.New("not implemented yet")
)

// Wallet struct handler
type Wallet struct {
	requestProvider provider.Provider
	isSigner        bool
	keyManager      *KeyManager
}

// NewRandomWallet returns Wallet from random values
func NewRandomWallet() (*Wallet, error) {
	keyManager, err := NewRandomKeyManager()
	if err != nil {
		return nil, err
	}

	return &Wallet{
		isSigner:   true,
		keyManager: keyManager,
	}, nil
}

// NewWalletFromPrivatekey returns Wallet from random values
func NewWalletFromPrivatekey(privateKey string) (*Wallet, error) {
	keyManager, err := NewKeyManagerFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		isSigner:   true,
		keyManager: keyManager,
	}, nil
}

// Connect assigns provider and changes isSigner accordingly
func (w *Wallet) Connect(requestProvider provider.Provider) {
	w.requestProvider = requestProvider
	w.isSigner = w.isSigner || w.keyManager.IsConnected()
}

// SignTransaction not implemented
// TODO: implement SignTransaction
func (w *Wallet) SignTransaction(req TransactionRequest) (string, error) {
	return "", ErrNotImplemented
}

// GetKeyManager returns KeyManager value
func (w *Wallet) GetKeyManager() *KeyManager {
	return w.keyManager
}

// GetProvider returns Provider value
func (w *Wallet) GetProvider() provider.Provider {
	return w.requestProvider
}

// IsSigner returns IsSigner value
func (w *Wallet) IsSigner() bool {
	return w.isSigner
}
