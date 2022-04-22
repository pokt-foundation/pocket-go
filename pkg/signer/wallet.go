package signer

// Wallet struct handler
type Wallet struct {
	keyManager *KeyManager
}

// NewRandomWallet returns Wallet from random values
func NewRandomWallet() (*Wallet, error) {
	keyManager, err := NewRandomKeyManager()
	if err != nil {
		return nil, err
	}

	return &Wallet{
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
		keyManager: keyManager,
	}, nil
}

// Sign returns a signed request as encoded hex string
func (w *Wallet) Sign(payload []byte) (string, error) {
	return w.keyManager.Sign(payload)
}

// SignBytes returns a signed request as raw bytes
func (w *Wallet) SignBytes(payload []byte) ([]byte, error) {
	return w.keyManager.SignBytes(payload)
}

// GetAddress returns address value
func (w *Wallet) GetAddress() string {
	return w.keyManager.address
}

// GetPublicKey returns public key value
func (w *Wallet) GetPublicKey() string {
	return w.keyManager.publicKey
}

// GetPrivateKey returns private key value
func (w *Wallet) GetPrivateKey() string {
	return w.keyManager.privateKey
}
