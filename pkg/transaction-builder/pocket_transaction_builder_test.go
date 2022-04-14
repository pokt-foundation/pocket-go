package transactionbuilder

import (
	"testing"

	"github.com/pokt-foundation/pocket-go/pkg/provider"
	"github.com/pokt-foundation/pocket-go/pkg/signer"
	"github.com/stretchr/testify/require"
)

func TestPocketTransactionBuilder_CreateTransaction(t *testing.T) {
	c := require.New(t)

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	msgSend, err := NewMsgSend("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561", 21)
	c.NoError(err)

	input, err := txBuilder.CreateTransaction("0021", "ohana", 21, msgSend, provider.Upokt)
	c.NotEmpty(input)
	c.NoError(err)
}
