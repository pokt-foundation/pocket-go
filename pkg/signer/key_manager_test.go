package signer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyManager_Sign(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	payload := "\"jsonrpc\":\"2.0\",\"method\":\"web3_sha3\",\"params\":[\"0x68656c6c6f20776f726c64\"],\"id\":64"
	expectedSignature := "3da5bdeaad60376b1b0f8d77bc33744bfee8c69f324057c84e83c79c1984bff7875d8a3c16e90d72eacb34da6e50d195852dcac0b74771799526e01e719b8502"

	keyManager, err := NewKeyManagerFromPrivateKey(privateKey)
	c.NoError(err)

	signature, err := keyManager.Sign([]byte(payload))
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
