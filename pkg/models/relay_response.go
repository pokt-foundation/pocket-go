package types

// RelayResponse interface that represents relay response
type RelayResponse interface {
	GetSignature() string
	GetPayload() string
	GetProof() RelayProof
	GetRelayRequest() RelayRequest
}
