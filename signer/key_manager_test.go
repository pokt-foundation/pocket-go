package signer

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyManager_Sign(t *testing.T) {
	c := require.New(t)

	keyManager, err := NewKeyManagerFromPrivateKey("")
	c.Equal(ErrMissingPrivateKey, err)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	payload := "deadbeef"
	expectedSignature := "5d04dfc0d0e579d815f761b452c7d01e5f20a71b9fb66dbbeb1959cffed9da0a621ee06dfd11171757f9c9541768eaf59cce75ac4acc1ad122556ec26e166108"

	keyManager, err = NewKeyManagerFromPrivateKey(privateKey)
	c.NoError(err)

	decodedPayload, err := hex.DecodeString(payload)
	c.NoError(err)

	signature, err := keyManager.Sign(decodedPayload)
	c.NoError(err)
	c.Equal(expectedSignature, signature)

	signatureBytes, err := keyManager.SignBytes(decodedPayload)
	c.NoError(err)
	c.Equal(expectedSignature, hex.EncodeToString(signatureBytes))
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
