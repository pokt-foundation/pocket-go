package signer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyManager_Sign(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	payload := "226a736f6e727063223a22322e30222c226d6574686f64223a22776562335f73686133222c22706172616d73223a5b22307836383635366336633666323037373666373236633634225d2c226964223a3634"
	expectedSignature := "3da5bdeaad60376b1b0f8d77bc33744bfee8c69f324057c84e83c79c1984bff7875d8a3c16e90d72eacb34da6e50d195852dcac0b74771799526e01e719b8502"

	keyManager, err := NewKeyManagerFromPrivateKey(privateKey)
	c.NoError(err)

	signature, err := keyManager.Sign(payload)
	c.NoError(err)
	c.Equal(expectedSignature, signature)
}

func TestKeyManager_GetAccount(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"

	keyManager, err := NewKeyManagerFromPrivateKey(privateKey)
	c.NoError(err)

	account := keyManager.GetAccount()
	c.Equal(keyManager.GetAddress(), account.Address)
	c.Equal(keyManager.GetPrivateKey(), account.PrivateKey)
	c.Equal(keyManager.GetPublicKey(), account.PublicKey)
}

func TestKeyManager_IsConnected(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"

	keyManager, err := NewKeyManagerFromPrivateKey(privateKey)
	c.NoError(err)
	c.True(keyManager.IsConnected())

	keyManager = &KeyManager{}
	c.False(keyManager.IsConnected())
}
