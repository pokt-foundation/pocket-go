package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAddress(t *testing.T) {
	c := require.New(t)

	c.True(ValidateAddress("1f32488b1db60fe528ab21e3cc26c96696be3faa"))
	c.False(ValidateAddress(";DROP DATABASE;"))
}

func TestValidatePrivateKey(t *testing.T) {
	c := require.New(t)

	c.True(ValidatePrivateKey("1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"))
	c.False(ValidatePrivateKey(";DROP DATABASE;"))
}

func TestValidatePublicKey(t *testing.T) {
	c := require.New(t)

	c.True(ValidatePublicKey("b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"))
	c.False(ValidatePublicKey(";DROP DATABASE;"))
}
