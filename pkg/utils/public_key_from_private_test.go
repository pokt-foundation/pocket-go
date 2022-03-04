package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicKeyFromPrivate(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	expectedPublicKey := "b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"

	c.Equal(expectedPublicKey, PublicKeyFromPrivate(privateKey))
}
