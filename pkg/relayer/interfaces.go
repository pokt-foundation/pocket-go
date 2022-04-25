package relayer

import "github.com/pokt-foundation/pocket-go/pkg/provider"

// Relayer interface that represents a relayer
type Relayer interface {
	Relay(input *RelayInput, options *provider.RelayRequestOptions) (*Output, error)
}

// Provider interface representing provider functions necessary for Relayer Package
type Provider interface {
	Relay(rpcURL string, input *provider.Relay, options *provider.RelayRequestOptions) (*provider.RelayOutput, error)
}

// Signer interface representing signer functions necessary for Relayer Package
type Signer interface {
	Sign(payload []byte) (string, error)
}
