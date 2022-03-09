package types

// RelayRequest interface represents a relay request
type RelayRequest interface {
	GetPayload() RelayPayload
	GetMeta() RelayMeta
	GetProof() RelayProof
}
