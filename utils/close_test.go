package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type closableMock struct {
	mock.Mock
}

func (c *closableMock) Close() error {
	args := c.Called()

	return args.Error(0)
}

func TestCloseOrLog(t *testing.T) {
	c := require.New(t)

	closableMock := &closableMock{}

	closableMock.On("Close").Return(nil).Once()

	c.NotPanics(func() {
		CloseOrLog(closableMock)
	})

	closableMock.On("Close").Return(errors.New("error")).Once()

	c.NotPanics(func() {
		CloseOrLog(closableMock)
	})
}
