package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewDefaultClient(t *testing.T) {
	c := require.New(t)

	client := NewDefaultClient()
	c.NotEmpty(client)

	client = NewCustomClient(5, 3*time.Second)
	c.NotEmpty(client)
}
