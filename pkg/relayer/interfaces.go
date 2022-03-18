package relayer

import "github.com/pokt-foundation/pocket-go/pkg/provider"

// Relayer interface that represents a relayer
type Relayer interface {
	GetNewSession(chain string, sessionHeight int, options *provider.DispatchRequestOptions) (*provider.Session, error)
	Relay(input *RelayInput, options *provider.RelayRequestOptions) (*provider.Relay, error)
}
