package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	address = "b50a6e20d3733fb89631ae32385b3c85c533c560"
)

func TestGetAddressFromPublickey(t *testing.T) {
	c := require.New(t)

	resultingAddress, err := GetAddressFromPublickey(publicKey)
	c.NoError(err)

	c.Equal(resultingAddress, address)
}
