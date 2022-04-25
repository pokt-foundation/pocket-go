package signer

// Signer interface that represents a Signer implementation
type Signer interface {
	Sign(payload []byte) (string, error)
	SignBytes(payload []byte) ([]byte, error)
	GetAddress() string
	GetPrivateKey() string
	GetPublicKey() string
}
