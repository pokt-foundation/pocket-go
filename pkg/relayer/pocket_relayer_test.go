package relayer

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/pocket-go/pkg/mock-client"
	"github.com/pokt-foundation/pocket-go/pkg/provider"
	"github.com/pokt-foundation/pocket-go/pkg/signer"
	"github.com/stretchr/testify/require"
)

func TestPocketRelayer_RelayerInterface(t *testing.T) {
	c := require.New(t)

	relayer := &PocketRelayer{}

	var i interface{} = relayer

	_, ok := i.(Relayer)
	c.True(ok)
}

func TestPocketRelayer_Relay(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	relayer := NewPocketRelayer(nil, nil)
	relayInput := &RelayInput{}

	relay, err := relayer.Relay(relayInput, nil)
	c.Equal(ErrNoSigner, err)
	c.Empty(relay)

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	relayer.signer = wallet

	relay, err = relayer.Relay(relayInput, nil)
	c.Equal(ErrNoProvider, err)
	c.Empty(relay)

	relayer.provider = provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"})

	relay, err = relayer.Relay(relayInput, nil)
	c.Equal(ErrNoSession, err)
	c.Empty(relay)

	relayInput.Session = &provider.Session{}

	relay, err = relayer.Relay(relayInput, nil)
	c.Equal(ErrNoPocketAAT, err)
	c.Empty(relay)

	relayInput.PocketAAT = &provider.PocketAAT{}

	relay, err = relayer.Relay(relayInput, nil)
	c.Equal(ErrSessionHasNoNodes, err)
	c.Empty(relay)

	relayInput.Session.Nodes = []*provider.Node{{PublicKey: "AOG"}}

	relay, err = relayer.Relay(relayInput, nil)
	c.Equal(ErrNoSessionHeader, err)
	c.Empty(relay)

	relayInput.Session.Header = &provider.SessionHeader{}
	relayInput.Node = &provider.Node{PublicKey: "PJOG"}

	relay, err = relayer.Relay(relayInput, nil)
	c.Equal(ErrNodeNotInSession, err)
	c.Empty(relay)

	relayInput.Node = nil

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusInternalServerError, "../provider/samples/client_relay.json")

	relay, err = relayer.Relay(relayInput, nil)
	c.Equal(provider.Err5xxOnConnection, err)
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusBadRequest, "../provider/samples/client_relay_error.json")

	var relayError *provider.RelayError

	relay, err = relayer.Relay(relayInput, nil)
	c.ErrorAs(err, &relayError)
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRelayRoute),
		http.StatusOK, "../provider/samples/client_relay.json")

	relay, err = relayer.Relay(relayInput, nil)
	c.NoError(err)
	c.NotEmpty(relay)

	relayInput.Node = &provider.Node{PublicKey: "AOG"}

	relay, err = relayer.Relay(relayInput, nil)
	c.NoError(err)
	c.NotEmpty(relay)
}
