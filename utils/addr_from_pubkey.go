package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetAddressFromDecodedPublickey converts an Application's Decoded key into an address
// decodedKey parameter is the application's public key decoded from hex string to []byte
// returns the application's address
func GetAddressFromDecodedPublickey(decodedKey []byte) (string, error) {
	hasher := sha256.New()

	_, err := hasher.Write(decodedKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil))[0:40], nil
}

// GetAddressFromPublickey converts an Application's Public key into an address
// publicKey parameter is the application's public key
// returns the application's address
func GetAddressFromPublickey(publicKey string) (string, error) {
	decodedKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()

	_, err = hasher.Write(decodedKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil))[0:40], nil
}
