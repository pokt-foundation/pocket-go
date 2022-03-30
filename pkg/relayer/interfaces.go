package relayer

import "github.com/pokt-foundation/pocket-go/pkg/provider"

// Relayer interface that represents a relayer
type Relayer interface {
	Relay(input *RelayInput, options *provider.RelayRequestOptions) (*Output, error)
}
