package utils

import "regexp"

const (
	addressLength    = 40
	publicKeyLength  = 64
	privateKeyLength = 128
)

var hexRegex = regexp.MustCompile("^[a-fA-F0-9]+$")

// ValidateAddress returns bool identifying if address is valid or not
func ValidateAddress(address string) bool {
	return len(address) == addressLength && hexRegex.MatchString(address)
}

// ValidatePrivateKey returns bool identifying if private key is valid or not
func ValidatePrivateKey(privateKey string) bool {
	return len(privateKey) == privateKeyLength && hexRegex.MatchString(privateKey)
}

// ValidatePublicKey returns bool identifying if public key is valid or not
func ValidatePublicKey(publicKey string) bool {
	return len(publicKey) == publicKeyLength && hexRegex.MatchString(publicKey)
}
