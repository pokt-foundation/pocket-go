package relayer

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/pocket-go/provider"
	"github.com/pokt-foundation/pocket-go/signer"
	"github.com/pokt-foundation/utils-go/mock-client"
	"github.com/stretchr/testify/require"
)

func TestRelayer_Relay(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	relayer := NewRelayer(nil, nil)
	input := &Input{}

	relay, _, err := relayer.Relay(input, nil)
	c.Equal(ErrNoSigner, err)
	c.Empty(relay)

	signer, signerErrr := signer.NewRandomSigner()
	c.NoError(signerErrr)

	relayer.signer = signer

	relay, _, err = relayer.Relay(input, nil)
	c.Equal(ErrNoProvider, err)
	c.Empty(relay)

	relayer.provider = provider.NewProvider("https://dummy.com", []string{"https://dummy.com"})

	relay, _, err = relayer.Relay(input, nil)
	c.Equal(ErrNoSession, err)
	c.Empty(relay)

	input.Session = &provider.Session{}

	relay, _, err = relayer.Relay(input, nil)
	c.Equal(ErrNoPocketAAT, err)
	c.Empty(relay)

	input.PocketAAT = &provider.PocketAAT{}

	relay, _, err = relayer.Relay(input, nil)
	c.Equal(ErrSessionHasNoNodes, err)
	c.Empty(relay)

	input.Session.Nodes = []provider.Node{{PublicKey: "AOG"}}

	relay, _, err = relayer.Relay(input, nil)
	c.Equal(ErrNoSessionHeader, err)
	c.Empty(relay)

	input.Session.Header = provider.SessionHeader{Chain: "chain"}
	input.Node = &provider.Node{PublicKey: "PJOG"}

	relay, _, err = relayer.Relay(input, nil)
	c.Equal(ErrNodeNotInSession, err)
	c.Empty(relay)

	input.Node = nil

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusInternalServerError, "../provider/samples/client_relay.json")

	relay, _, err = relayer.Relay(input, nil)
	c.NotNil(err)
	c.Equal(provider.Err5xxOnConnection, err)
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusBadRequest, "../provider/samples/client_relay_error.json")

	var error *provider.RelayError

	relay, _, err = relayer.Relay(input, nil)
	c.ErrorAs(err, &error)
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusOK, "../provider/samples/client_relay.json")

	relay, _, err = relayer.Relay(input, nil)
	c.Nil(err)
	c.NotEmpty(relay)

	input.Node = &provider.Node{PublicKey: "AOG"}

	relay, _, err = relayer.Relay(input, nil)
	c.Nil(err)
	c.NotEmpty(relay)
}

func TestRelayer_RelayWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	relayer := NewRelayer(nil, nil)
	input := &Input{}

	relay, _, err := relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(ErrNoSigner, err)
	c.Empty(relay)

	signer, signerErr := signer.NewRandomSigner()
	c.NoError(signerErr)

	relayer.signer = signer

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(ErrNoProvider, err)
	c.Empty(relay)

	relayer.provider = provider.NewProvider("https://dummy.com", []string{"https://dummy.com"})

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(ErrNoSession, err)
	c.Empty(relay)

	input.Session = &provider.Session{}

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(ErrNoPocketAAT, err)
	c.Empty(relay)

	input.PocketAAT = &provider.PocketAAT{}

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(ErrSessionHasNoNodes, err)
	c.Empty(relay)

	input.Session.Nodes = []provider.Node{{PublicKey: "AOG"}}

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(ErrNoSessionHeader, err)
	c.Empty(relay)

	input.Session.Header = provider.SessionHeader{Chain: "chain"}
	input.Node = &provider.Node{PublicKey: "PJOG"}

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(ErrNodeNotInSession, err)
	c.Empty(relay)

	input.Node = nil

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusInternalServerError, "../provider/samples/client_relay.json")

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Equal(provider.Err5xxOnConnection, err)
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusBadRequest, "../provider/samples/client_relay_error.json")

	var error *provider.RelayError

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.ErrorAs(err, &error)
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusOK, "../provider/samples/client_relay.json")

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Nil(err)
	c.NotEmpty(relay)

	input.Node = &provider.Node{PublicKey: "AOG"}

	relay, _, err = relayer.RelayWithCtx(context.Background(), input, nil)
	c.Nil(err)
	c.NotEmpty(relay)
}
