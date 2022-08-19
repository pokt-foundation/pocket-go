// Package signer creates crypto safe signature with different inputs
package signer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"regexp"
	"strconv"

	"github.com/pokt-foundation/pocket-go/utils"
	"golang.org/x/crypto/scrypt"
)

const (
	scryptHashLength = 32
	scryptN          = 32768
	scryptR          = 8
	scryptP          = 1
	scryptSec        = 12
	tagLength        = 16
	defaultKDF       = "scrypt"
)

var (
	// ErrInvalidPrivateKey error when private key is invalid
	ErrInvalidPrivateKey = errors.New("invalid private key")
	// ErrInvalidPPK error when PPK is invalid
	ErrInvalidPPK = errors.New("invalid ppk")

	base64Regex = regexp.MustCompile("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$")
	hexRegex    = regexp.MustCompile("^[a-fA-F0-9]+$")
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
	if !utils.ValidatePrivateKey(privateKey) {
		return nil, ErrInvalidPrivateKey
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

func getAESGCMValues(password, saltBytes []byte) ([]byte, cipher.AEAD, error) {
	scryptKey, err := scrypt.Key(password, saltBytes, scryptN, scryptR, scryptP, scryptHashLength)
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(scryptKey)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	return scryptKey[:scryptSec], gcm, nil
}

// NewSignerFromPPK returns Signer from PPK and its password
func NewSignerFromPPK(password string, ppk *PPK) (*Signer, error) {
	if !ppk.Validate() {
		return nil, ErrInvalidPPK
	}

	saltBytes, err := hex.DecodeString(ppk.Salt)
	if err != nil {
		return nil, err
	}

	encBytes, err := base64.StdEncoding.DecodeString(ppk.Ciphertext)
	if err != nil {
		return nil, err
	}

	nuance, gcm, err := getAESGCMValues([]byte(password), saltBytes)
	if err != nil {
		return nil, err
	}

	privateKey, err := gcm.Open(nil, nuance, encBytes, nil)
	if err != nil {
		return nil, err
	}

	return NewSignerFromPrivateKey(string(privateKey))
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

// PPK struct handler for Portable Private Key file
// Do not change this struct to stay compatible
type PPK struct {
	Kdf        string `json:"kdf"`
	Salt       string `json:"salt"`
	SecParam   string `json:"secparam"`
	Hint       string `json:"hint"`
	Ciphertext string `json:"ciphertext"`
}

func randBytes(numBytes int) ([]byte, error) {
	b := make([]byte, numBytes)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// NewPPK returns an instance of PPK for given input
func NewPPK(privateKey, password, hint string) (*PPK, error) {
	saltBytes, err := randBytes(tagLength)
	if err != nil {
		return nil, err
	}

	nuance, gcm, err := getAESGCMValues([]byte(password), saltBytes)
	if err != nil {
		return nil, err
	}

	return &PPK{
		Kdf:        defaultKDF,
		Salt:       hex.EncodeToString(saltBytes),
		SecParam:   strconv.Itoa(scryptSec),
		Hint:       hint,
		Ciphertext: base64.StdEncoding.EncodeToString(gcm.Seal(nil, nuance, []byte(privateKey), nil)),
	}, nil
}

// Validate returns bool representing if PPK instance is valid or not
func (ppk *PPK) Validate() bool {
	secParam, err := strconv.Atoi(ppk.SecParam)
	if err != nil {
		return false
	}

	return ppk.Kdf == defaultKDF &&
		ppk.Salt != "" &&
		hexRegex.MatchString(ppk.Salt) &&
		secParam >= 0 &&
		ppk.Ciphertext != "" &&
		base64Regex.MatchString(ppk.Ciphertext)
}
