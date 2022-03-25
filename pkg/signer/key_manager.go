package signer

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"

	"github.com/pokt-foundation/pocket-go/pkg/utils"
)

var (
	// ErrCryptoSignDetached is error when ErrCryptoSignDetached function exits value other than 0
	ErrCryptoSignDetached = errors.New("error in CryptoSignDetached")
	// ErrCryptoSignKeyPair is error when ErrCryptoSignKeyPair function exits value other than 0
	ErrCryptoSignKeyPair = errors.New("error in CryptoSignKeyPair")
)

// KeyManager struct handler
type KeyManager struct {
	address    string
	publicKey  string
	privateKey string
}

// NewRandomKeyManager returns a KeyManager with random keys
func NewRandomKeyManager() (*KeyManager, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	address, err := utils.GetAddressFromDecodedPublickey(publicKey)
	if err != nil {
		return nil, err
	}

	return &KeyManager{
		address:    address,
		publicKey:  hex.EncodeToString(publicKey),
		privateKey: hex.EncodeToString(privateKey),
	}, nil
}

// NewKeyManagerFromPrivateKey returns KeyManager from private key
func NewKeyManagerFromPrivateKey(privateKey string) (*KeyManager, error) {
	publicKey := utils.PublicKeyFromPrivate(privateKey)

	address, err := utils.GetAddressFromPublickey(publicKey)
	if err != nil {
		return nil, err
	}

	return &KeyManager{
		address:    address,
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

// Sign returns a signed request
// Note that libsodium uses detached signatures
func (km *KeyManager) Sign(payload []byte) (string, error) {
	decodedKey, err := hex.DecodeString(km.privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ed25519.Sign(decodedKey, payload)), nil
}

// GetAddress returns address value
func (km *KeyManager) GetAddress() string {
	return km.address
}

// GetPublicKey returns public key value
func (km *KeyManager) GetPublicKey() string {
	return km.publicKey
}

// GetPrivateKey returns private key value
func (km *KeyManager) GetPrivateKey() string {
	return km.privateKey
}

// GetAccount returns Account struct holding all key values
func (km *KeyManager) GetAccount() *Account {
	return &Account{
		Address:    km.address,
		PublicKey:  km.publicKey,
		PrivateKey: km.privateKey,
	}
}

// IsConnected returns if KeyManger is connected or not
func (km *KeyManager) IsConnected() bool {
	return km.address != "" && km.publicKey != "" && km.privateKey != ""
}
