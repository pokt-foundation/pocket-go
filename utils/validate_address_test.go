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
