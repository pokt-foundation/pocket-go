package relayer

import (
	"github.com/pokt-foundation/pocket-go/pkg/models"
	"github.com/pokt-foundation/pocket-go/pkg/provider"
)

// RelayInput struct that represents data needed for doing a relay request
type RelayInput struct {
	Blockchain string
	Data       string
	Headers    models.RelayHeaders
	Method     string
	Node       *provider.Node
	Path       string
	PocketAAT  *provider.PocketAAT
	Session    *provider.Session
}

// RequestHash struct holding data needed to create a request hash
type RequestHash struct {
	Payload *provider.RelayPayload `json:"payload"`
	Meta    *provider.RelayMeta    `json:"meta"`
}

// RelayerResponse struct for data needed as response for relay request
type RelayerResponse struct {
	Response *provider.RelayResponse
	Proof    *provider.RelayProof
	Node     *provider.Node
}

// Order of fields matters for signature
type relayProofForSignature struct {
	Entropy            int    `json:"entropy"`
	SessionBlockHeight int    `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Signature          string `json:"signature"`
	Token              string `json:"token"`
	RequestHash        string `json:"request_hash"`
}
