package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetAddressFromPublickey converts an Application's Public Key into an address
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
