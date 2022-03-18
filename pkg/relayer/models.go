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
