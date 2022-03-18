package provider

import (
	"fmt"

	"github.com/pokt-foundation/pocket-go/pkg/models"
)

// RelayResponse struct wrapper for both possible responses of a relay request
type RelayResponse struct {
	SuccessfulResponse *Relay
	ErrorResponse      *RelayErrorResponse
}

// Relay represents a Relay to Pocket
type Relay struct {
	Payload *RelayPayload `json:"payload"`
	Meta    *RelayMeta    `json:"meta"`
	Proof   *RelayProof   `json:"proof"`
}

// RelayMeta represents metadata of a relay
type RelayMeta struct {
	BlockHeight int `json:"block_height"`
}

// RelayPayload represents payload of a relay
type RelayPayload struct {
	Data    string              `json:"data"`
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers models.RelayHeaders `json:"headers"`
}

// RelayProof represents proof of a relay
type RelayProof struct {
	RequestHash        string     `json:"request_hash"`
	Entropy            int64      `json:"entropy"`
	SessionBlockHeight int        `json:"session_block_height"`
	ServicerPubKey     string     `json:"servicer_pub_key"`
	Blockchain         string     `json:"blockchain"`
	AAT                *PocketAAT `json:"aat"`
	Signature          string     `json:"signature"`
}

// PocketAAT represents a Pocket Application Authentication Token
type PocketAAT struct {
	Version      string `json:"version"`
	AppPubKey    string `json:"app_pub_key"`
	ClientPubKey string `json:"client_pub_key"`
	Signature    string `json:"signature"`
}

// RelayErrorResponse represents error response of relay request
type RelayErrorResponse struct {
	Error    *RelayError       `json:"error"`
	Dispatch *DispatchResponse `json:"dispatch"`
}

// RelayError represents the thrown error of a relay request
type RelayError struct {
	Code      int    `json:"code"`
	Codespace string `json:"codespace"`
	Message   string `json:"message"`
}

// Error returns string representation of error
// needed to implement error interface
func (e *RelayError) Error() string {
	return fmt.Sprintf("Request failed with code: %v, codespace: %s and message: %s", e.Code, e.Codespace, e.Message)
}
