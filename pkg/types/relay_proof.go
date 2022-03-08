package types

// RelayProof interface represents a relay proof
type RelayProof interface {
	GetEntropy() string
	GetSessionBlockHeight() int
	GetServicerPubKey() string
	GetBlockchain() string
	GetToken() PocketAAT
	GetSignature() string
	GetRequestHash() string
}
