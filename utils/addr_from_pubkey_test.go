package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAddressFromPublickey(t *testing.T) {
	c := require.New(t)

	publicKey := "b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	expectedAddress := "b50a6e20d3733fb89631ae32385b3c85c533c560"

	resultingAddress, err := GetAddressFromPublickey(publicKey)
	c.NoError(err)

	c.Equal(resultingAddress, expectedAddress)
}

func TestGetAddressFromDecodedPublickey(t *testing.T) {
	c := require.New(t)

	publicKey := "b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	expectedAddress := "b50a6e20d3733fb89631ae32385b3c85c533c560"

	decodedKey, err := hex.DecodeString(publicKey)
	c.NoError(err)

	resultingAddress, err := GetAddressFromDecodedPublickey(decodedKey)
	c.NoError(err)

	c.Equal(resultingAddress, expectedAddress)
}
