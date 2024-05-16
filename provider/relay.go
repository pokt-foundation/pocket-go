package provider

import (
	"fmt"
)

// RelayInput represents input needed to do a Relay to Pocket
type RelayInput struct {
	Payload *RelayPayload `json:"payload"`
	Meta    *RelayMeta    `json:"meta"`
	Proof   *RelayProof   `json:"proof"`
}

// RelayOutput represents the Relay RPC output
type RelayOutput struct {
	Response  string `json:"response"`
	Signature string `json:"signature"`
}

// RelayMeta represents metadata of a relay
type RelayMeta struct {
	BlockHeight int `json:"block_height"`
}

// RelayHeaders map of relay headers
type RelayHeaders map[string]string

// RelayPayload represents payload of a relay
type RelayPayload struct {
	Data    string       `json:"data"`
	Method  string       `json:"method"`
	Path    string       `json:"path"`
	Headers RelayHeaders `json:"headers"`
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

// RelayErrorOutput represents error response of relay request
type RelayErrorOutput struct {
	Error struct {
		Code      RelayErrorCode `json:"code"`
		Codespace string         `json:"codespace"`
		Message   string         `json:"message"`
	} `json:"error"`
}

// RelayError represents the thrown error of a relay request
type RelayError struct {
	Code           RelayErrorCode
	Codespace      string
	Message        string
	ServicerPubKey string
}

// Error returns string representation of error
// needed to implement error interface
func (e *RelayError) Error() string {
	return fmt.Sprintf("Request failed with code: %v, codespace: %s and message: %s\nWith ServicerPubKey: %s",
		e.Code, e.Codespace, e.Message, e.ServicerPubKey)
}

// RelayErrorCode is enum of possible relay error codes
type RelayErrorCode int

const (
	// AppNotFoundError error when app is not found
	AppNotFoundError RelayErrorCode = 45
	// DuplicateProofError error when proof is used multiple times
	DuplicateProofError RelayErrorCode = 37
	// EmptyPayloadDataError error when sent payload is empty
	EmptyPayloadDataError RelayErrorCode = 25
	// EvidencedSealedError error when evidence is sealed, either max relays reached or claim already submitted
	EvidencedSealedError RelayErrorCode = 90
	// HTTPExecutionError error when http request failed
	HTTPExecutionError RelayErrorCode = 28
	// InvalidBlockHeightError error when sent block height is invalid
	InvalidBlockHeightError RelayErrorCode = 60
	// InvalidSessionError error when session is invalid
	InvalidSessionError RelayErrorCode = 14
	// OutOfSyncRequestError error when request is not on sync
	OutOfSyncRequestError RelayErrorCode = 75
	// OverServiceError error when request exceeds service capacity
	OverServiceError RelayErrorCode = 71
	// RequestHashError error when hash is not correct
	RequestHashError RelayErrorCode = 74
	// UnsupportedBlockchainError error when sent blockchain is not supported yet
	UnsupportedBlockchainError RelayErrorCode = 76
)

// IsErrorCode returns if error has the same relay error code as input
func IsErrorCode(code RelayErrorCode, err error) bool {
	castedErr, ok := err.(*RelayError)
	if !ok {
		return false
	}

	return castedErr.Code == code
}
