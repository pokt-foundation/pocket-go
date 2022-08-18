package signer

import (
	"encoding/hex"
	"testing"

	"github.com/pokt-foundation/pocket-go/utils"
	"github.com/stretchr/testify/require"
)

func TestNewRandomSigner(t *testing.T) {
	c := require.New(t)

	signer, err := NewRandomSigner()
	c.NoError(err)

	c.Equal(utils.PublicKeyFromPrivate(signer.GetPrivateKey()), signer.GetPublicKey())

	expectedAddres, err := utils.GetAddressFromPublickey(signer.GetPublicKey())
	c.NoError(err)
	c.Equal(expectedAddres, signer.GetAddress())
}

func TestNewSignerFromPPK(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	password := "bebitofiufiu"
	hint := "fiufiu"

	ppk, err := NewPPK(privateKey, password, hint)
	c.NoError(err)
	c.True(ppk.Validate())

	ppkSigner, err := NewSignerFromPPK(password, &PPK{SecParam: "FIU"})
	c.Equal(ErrInvalidPPK, err)
	c.Empty(ppkSigner)

	ppkSigner, err = NewSignerFromPPK(password, ppk)
	c.NoError(err)
	c.NotEmpty(ppkSigner)

	fromPrivKeySigner, err := NewSignerFromPrivateKey(privateKey)
	c.NoError(err)
	c.NotEmpty(fromPrivKeySigner)

	c.Equal(fromPrivKeySigner.GetAddress(), ppkSigner.GetAddress())
	c.Equal(fromPrivKeySigner.GetPrivateKey(), ppkSigner.GetPrivateKey())
	c.Equal(fromPrivKeySigner.GetPublicKey(), ppkSigner.GetPublicKey())
}

func TestSigner_Sign(t *testing.T) {
	c := require.New(t)

	signer, err := NewSignerFromPrivateKey("")
	c.Equal(ErrInvalidPrivateKey, err)
	c.Empty(signer)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"
	payload := "deadbeef"
	expectedSignature := "5d04dfc0d0e579d815f761b452c7d01e5f20a71b9fb66dbbeb1959cffed9da0a621ee06dfd11171757f9c9541768eaf59cce75ac4acc1ad122556ec26e166108"

	signer, err = NewSignerFromPrivateKey(privateKey)
	c.NoError(err)

	decodedPayload, err := hex.DecodeString(payload)
	c.NoError(err)

	signature, err := signer.Sign(decodedPayload)
	c.NoError(err)
	c.Equal(expectedSignature, signature)

	signatureBytes, err := signer.SignBytes(decodedPayload)
	c.NoError(err)
	c.Equal(expectedSignature, hex.EncodeToString(signatureBytes))
}

func TestSigner_GetAccount(t *testing.T) {
	c := require.New(t)

	privateKey := "1f8cbde30ef5a9db0a5a9d5eb40536fc9defc318b8581d543808b7504e0902bcb243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3"

	signer, err := NewSignerFromPrivateKey(privateKey)
	c.NoError(err)

	account := signer.GetAccount()
	c.Equal(signer.GetAddress(), account.Address)
	c.Equal(signer.GetPrivateKey(), account.PrivateKey)
	c.Equal(signer.GetPublicKey(), account.PublicKey)
}
