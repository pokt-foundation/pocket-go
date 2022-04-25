package utils

// PublicKeyFromPrivate extracts the public key from a 64-byte long ed25519 private key
// privateKey parameter is the private key buffer
// returns public key buffer
func PublicKeyFromPrivate(privateKey string) string {
	return privateKey[64:]
}
