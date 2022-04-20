package transactionbuilder

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/pocket-go/pkg/mock-client"
	"github.com/pokt-foundation/pocket-go/pkg/provider"
	"github.com/pokt-foundation/pocket-go/pkg/signer"
	"github.com/stretchr/testify/require"
)

func TestPocketTransactionBuilder_TransactionBuilderInterface(t *testing.T) {
	c := require.New(t)

	transactionbuilder := &PocketTransactionBuilder{}

	var i interface{} = transactionbuilder

	_, ok := i.(TransactionBuilder)
	c.True(ok)
}

func TestPocketTransactionBuilder_SubmitMsgSend(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	msgSend, err := NewMsgSend("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561", 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit("0021", "ohana", 21, msgSend, Upokt)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit("0021", "ohana", 21, msgSend, Upokt)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestPocketTransactionBuilder_SubmitStakeApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	stakeApp, err := NewStakeApp("b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3", []string{"0021"}, 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit("0021", "", 21, stakeApp, Upokt)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit("0021", "ohana", 21, stakeApp, Upokt)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestPocketTransactionBuilder_SubmitUnstakeApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	unstakeApp, err := NewUnstakeApp("b50a6e20d3733fb89631ae32385b3c85c533c560")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit("0021", "", 21, unstakeApp, Upokt)
	c.NoError(err)
	c.NotEmpty(output)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit("0021", "ohana", 21, unstakeApp, Upokt)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestPocketTransactionBuilder_SubmitUnjailApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	unjailApp, err := NewUnjailApp("b50a6e20d3733fb89631ae32385b3c85c533c560")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit("0021", "", 21, unjailApp, Upokt)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit("0021", "ohana", 21, unjailApp, Upokt)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestPocketTransactionBuilder_SubmitStakeNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	stakeNode, err := NewStakeNode("b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3", "https://dummy.com:443", "b50a6e20d3733fb89631ae32385b3c85c533c560", []string{"0021"}, 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit("0021", "", 21, stakeNode, Upokt)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit("0021", "ohana", 21, stakeNode, Upokt)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestPocketTransactionBuilder_SubmitUnstakeNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	unstakeNode, err := NewUnstakeNode("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit("0021", "", 21, unstakeNode, Upokt)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit("0021", "ohana", 21, unstakeNode, Upokt)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestPocketTransactionBuilder_SubmitUnjailNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	wallet, err := signer.NewRandomWallet()
	c.NoError(err)

	txBuilder := NewPocketTransactionBuilder(provider.NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}), wallet)

	unjailNode, err := NewUnjailNode("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit("0021", "", 21, unjailNode, Upokt)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit("0021", "ohana", 21, unjailNode, Upokt)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}
