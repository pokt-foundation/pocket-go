package signer

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"

	"github.com/pokt-foundation/pocket-go/utils"
)

var (
	// ErrMissingPrivateKey error when private key is not sent
	ErrMissingPrivateKey = errors.New("missing private key")
)

// Signer struct handler
type Signer struct {
	address    string
	publicKey  string
	privateKey string
}

// NewRandomSigner returns a Signer with random keys
func NewRandomSigner() (*Signer, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	address, err := utils.GetAddressFromDecodedPublickey(publicKey)
	if err != nil {
		return nil, err
	}

	return &Signer{
		address:    address,
		publicKey:  hex.EncodeToString(publicKey),
		privateKey: hex.EncodeToString(privateKey),
	}, nil
}

// NewSignerFromPrivateKey returns Signer from private key
func NewSignerFromPrivateKey(privateKey string) (*Signer, error) {
	if privateKey == "" {
		return nil, ErrMissingPrivateKey
	}

	publicKey := utils.PublicKeyFromPrivate(privateKey)

	address, err := utils.GetAddressFromPublickey(publicKey)
	if err != nil {
		return nil, err
	}

	return &Signer{
		address:    address,
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

// Sign returns a signed request as encoded hex string
func (s *Signer) Sign(payload []byte) (string, error) {
	decodedKey, err := hex.DecodeString(s.privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ed25519.Sign(decodedKey, payload)), nil
}

// SignBytes returns a signed request as raw bytes
func (s *Signer) SignBytes(payload []byte) ([]byte, error) {
	decodedKey, err := hex.DecodeString(s.privateKey)
	if err != nil {
		return nil, err
	}

	return ed25519.Sign(decodedKey, payload), nil
}

// GetAddress returns address value
func (s *Signer) GetAddress() string {
	return s.address
}

// GetPublicKey returns public key value
func (s *Signer) GetPublicKey() string {
	return s.publicKey
}

// GetPrivateKey returns private key value
func (s *Signer) GetPrivateKey() string {
	return s.privateKey
}

// Account holds an account's data
type Account struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

// GetAccount returns Account struct holding all key values
func (s *Signer) GetAccount() *Account {
	return &Account{
		Address:    s.address,
		PublicKey:  s.publicKey,
		PrivateKey: s.privateKey,
	}
}
