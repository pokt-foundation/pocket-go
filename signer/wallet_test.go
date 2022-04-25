package signer

import (
	"encoding/hex"
	"testing"

	"github.com/pokt-foundation/pocket-go/utils"
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

	c.Equal(utils.PublicKeyFromPrivate(wallet.GetPrivateKey()), wallet.GetPublicKey())

	expectedAddres, err := utils.GetAddressFromPublickey(wallet.GetPublicKey())
	c.NoError(err)
	c.Equal(expectedAddres, wallet.GetAddress())
}

func TestNewWalletFromPrivatekey(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"

	wallet, err := NewWalletFromPrivatekey(privateKey)
	c.NoError(err)

	c.Equal(utils.PublicKeyFromPrivate(wallet.GetPrivateKey()), wallet.GetPublicKey())

	expectedAddres, err := utils.GetAddressFromPublickey(wallet.GetPublicKey())
	c.NoError(err)
	c.Equal(expectedAddres, wallet.GetAddress())
}

func TestWallet_Sign(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	payload := "deadbeef"
	expectedSignature := "5d04dfc0d0e579d815f761b452c7d01e5f20a71b9fb66dbbeb1959cffed9da0a621ee06dfd11171757f9c9541768eaf59cce75ac4acc1ad122556ec26e166108"

	wallet, err := NewWalletFromPrivatekey(privateKey)
	c.NoError(err)

	decodedPayload, err := hex.DecodeString(payload)
	c.NoError(err)

	signature, err := wallet.Sign(decodedPayload)
	c.NoError(err)
	c.Equal(expectedSignature, signature)

	signatureBytes, err := wallet.SignBytes(decodedPayload)
	c.NoError(err)
	c.Equal(expectedSignature, hex.EncodeToString(signatureBytes))
}
