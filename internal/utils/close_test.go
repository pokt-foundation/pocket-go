package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type dummyClosable struct{}

func (c dummyClosable) Close() error { return nil }

func TestCloseOrLog(t *testing.T) {
	c := require.New(t)

	c.NotPanics(func() {
		CloseOrLog(dummyClosable{})
	})
}
