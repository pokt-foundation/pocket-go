package mock

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/pocket-go/internal/client"
	"github.com/stretchr/testify/require"
)

func TestAddMockedResponseFromFile(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	AddMockedResponseFromFile(http.MethodGet, "https://dummy.com", http.StatusCreated, "samples/dummy.json")

	client := client.NewDefaultClient()

	response, err := client.Get("https://dummy.com", http.Header{})
	c.Nil(err)
	c.NotNil(response)
	c.Equal(http.StatusCreated, response.StatusCode)
	c.NoError(response.Body.Close())

	c.Panics(func() {
		AddMockedResponseFromFile(http.MethodGet, "https://dummy.com", http.StatusCreated, "samples/not_found.json")
	})
}
