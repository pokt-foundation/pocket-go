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

func TestAddMultipleMockedResponses(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	AddMultipleMockedResponses(http.MethodGet, "https://dummy.com", http.StatusOK, []string{
		"samples/dummy.json",
		"samples/dummy.json",
	})

	client := client.NewDefaultClient()

	response1, err := client.Get("https://dummy.com", http.Header{})
	c.Nil(err)
	c.NotNil(response1)
	c.Equal(http.StatusOK, response1.StatusCode)
	c.NoError(response1.Body.Close())

	response2, err := client.Get("https://dummy.com", http.Header{})
	c.Nil(err)
	c.NotNil(response2)
	c.Equal(http.StatusOK, response2.StatusCode)
	c.NoError(response2.Body.Close())

	response3, err := client.Get("https://dummy.com", http.Header{})
	c.Nil(response3)
	c.Error(ErrResponseNotFound, err)
}
