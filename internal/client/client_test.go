package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/pocket-go/internal/mock-client"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultClient(t *testing.T) {
	c := require.New(t)

	client := NewDefaultClient()
	c.NotEmpty(client)

	client = NewCustomClient(5, 3*time.Second)
	c.NotEmpty(client)
}

func TestClient_PostWithURLJSONParams(t *testing.T) {
	c := require.New(t)

	client := NewDefaultClient()
	c.NotEmpty(client)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock.AddMockedResponseFromFile(http.MethodPost, "https://dummy.com", http.StatusCreated, "../mock-client/samples/dummy.json")

	response, err := client.PostWithURLJSONParams("https://dummy.com", map[string]string{
		"ohana": "family",
	}, http.Header{})
	c.NoError(err)

	c.NotEmpty(response)
	c.Equal(http.StatusCreated, response.StatusCode)
	c.NoError(response.Body.Close())
}
