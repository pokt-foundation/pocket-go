package signer

import (
	"testing"

	"github.com/pokt-foundation/pocket-go/pkg/provider"
	"github.com/pokt-foundation/pocket-go/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestWallet_SignerInterface(t *testing.T) {
	c := require.New(t)

	wallet := &Wallet{}

	var i interface{} = wallet

	_, ok := i.(Signer)
	c.True(ok)
}

func TestNewRandomWallet(t *testing.T) {
	c := require.New(t)

	wallet, err := NewRandomWallet()
	c.NoError(err)

	c.True(wallet.IsSigner())
	c.Equal(utils.PublicKeyFromPrivate(wallet.GetKeyManager().GetPrivateKey()), wallet.GetKeyManager().GetPublicKey())

	expectedAddres, err := utils.GetAddressFromPublickey(wallet.GetKeyManager().GetPublicKey())
	c.NoError(err)
	c.Equal(expectedAddres, wallet.GetKeyManager().GetAddress())
}

func TestNewWalletFromPrivatekey(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"

	wallet, err := NewWalletFromPrivatekey(privateKey)
	c.NoError(err)

	c.True(wallet.IsSigner())
	c.Equal(utils.PublicKeyFromPrivate(wallet.GetKeyManager().GetPrivateKey()), wallet.GetKeyManager().GetPublicKey())

	expectedAddres, err := utils.GetAddressFromPublickey(wallet.GetKeyManager().GetPublicKey())
	c.NoError(err)
	c.Equal(expectedAddres, wallet.GetKeyManager().GetAddress())
}

func TestWallet_Connect(t *testing.T) {
	c := require.New(t)

	wallet, err := NewRandomWallet()
	c.NoError(err)

	wallet.Connect(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}))

	c.True(wallet.IsSigner())
	c.NotEmpty(wallet.GetProvider())
}

func TestWallet_SignTransaction(t *testing.T) {
	c := require.New(t)

	wallet, err := NewRandomWallet()
	c.NoError(err)

	_, err = wallet.SignTransaction(nil)
	c.Equal(ErrNotImplemented, err)
}
