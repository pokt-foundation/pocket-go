package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type dummyClosable struct{}

func (c dummyClosable) Close() error { return nil }

type dummyErrorClosable struct{}

func (c dummyErrorClosable) Close() error { return errors.New("error") }

func TestCloseOrLog(t *testing.T) {
	c := require.New(t)

	c.NotPanics(func() {
		CloseOrLog(dummyClosable{})
	})

	c.NotPanics(func() {
		CloseOrLog(dummyErrorClosable{})
	})
}
