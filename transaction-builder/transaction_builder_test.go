package transactionbuilder

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/pocket-go/provider"
	"github.com/pokt-foundation/pocket-go/signer"
	"github.com/pokt-foundation/utils-go/mock-client"
	"github.com/stretchr/testify/require"
)

func TestTransactionBuilder_SubmitError(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	txBuilder := NewTransactionBuilder(nil, nil)

	output, err := txBuilder.Submit("", nil, nil)
	c.Empty(output)
	c.Equal(ErrNoProvider, err)

	txBuilder.provider = provider.NewProvider("https://dummy.com", []string{"https://dummy.com"})

	output, err = txBuilder.Submit("", nil, nil)
	c.Empty(output)
	c.Equal(ErrNoSigner, err)

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder.signer = signer

	output, err = txBuilder.Submit("", nil, nil)
	c.Empty(output)
	c.Equal(ErrNoChainID, err)

	output, err = txBuilder.Submit(Mainnet, nil, nil)
	c.Empty(output)
	c.Equal(ErrNoTransactionMessage, err)
}

func TestTransactionBuilder_SubmitMsgSend(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	msgSend, err := NewSend("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561", 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit(Mainnet, msgSend, &TransactionOptions{
		CoinDenom: Upokt,
		Fee:       23,
		Memo:      "ohana",
	})
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit(Mainnet, msgSend, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitStakeApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	stakeApp, err := NewStakeApp("b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3", []string{"0021"}, 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit(Mainnet, stakeApp, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit(Mainnet, stakeApp, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnstakeApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unstakeApp, err := NewUnstakeApp("b50a6e20d3733fb89631ae32385b3c85c533c560")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit(Mainnet, unstakeApp, nil)
	c.NoError(err)
	c.NotEmpty(output)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit(Mainnet, unstakeApp, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnjailApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unjailApp, err := NewUnjailApp("b50a6e20d3733fb89631ae32385b3c85c533c560")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit(Mainnet, unjailApp, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit(Mainnet, unjailApp, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitStakeNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	stakeNode, err := NewStakeNode("b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3", "https://dummy.com:443",
		"b50a6e20d3733fb89631ae32385b3c85c533c560", []string{"0021"}, 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit(Mainnet, stakeNode, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit(Mainnet, stakeNode, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnstakeNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unstakeNode, err := NewUnstakeNode("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit(Mainnet, unstakeNode, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit(Mainnet, unstakeNode, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnjailNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unjailNode, err := NewUnjailNode("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.Submit(Mainnet, unjailNode, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.Submit(Mainnet, unjailNode, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitErrorWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	txBuilder := NewTransactionBuilder(nil, nil)

	output, err := txBuilder.SubmitWithCtx(context.Background(), "", nil, nil)
	c.Empty(output)
	c.Equal(ErrNoProvider, err)

	txBuilder.provider = provider.NewProvider("https://dummy.com", []string{"https://dummy.com"})

	output, err = txBuilder.SubmitWithCtx(context.Background(), "", nil, nil)
	c.Empty(output)
	c.Equal(ErrNoSigner, err)

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder.signer = signer

	output, err = txBuilder.SubmitWithCtx(context.Background(), "", nil, nil)
	c.Empty(output)
	c.Equal(ErrNoChainID, err)

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, nil, nil)
	c.Empty(output)
	c.Equal(ErrNoTransactionMessage, err)
}

func TestTransactionBuilder_SubmitMsgSendWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	msgSend, err := NewSend("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561", 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.SubmitWithCtx(context.Background(), Mainnet, msgSend, &TransactionOptions{
		CoinDenom: Upokt,
		Fee:       23,
		Memo:      "ohana",
	})
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, msgSend, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitStakeAppWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	stakeApp, err := NewStakeApp("b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3", []string{"0021"}, 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.SubmitWithCtx(context.Background(), Mainnet, stakeApp, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, stakeApp, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnstakeAppWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unstakeApp, err := NewUnstakeApp("b50a6e20d3733fb89631ae32385b3c85c533c560")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.SubmitWithCtx(context.Background(), Mainnet, unstakeApp, nil)
	c.NoError(err)
	c.NotEmpty(output)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, unstakeApp, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnjailAppWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unjailApp, err := NewUnjailApp("b50a6e20d3733fb89631ae32385b3c85c533c560")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.SubmitWithCtx(context.Background(), Mainnet, unjailApp, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, unjailApp, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitStakeNodeWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	stakeNode, err := NewStakeNode("b243b27bc9fbe5580457a46370ae5f03a6f6753633e51efdaf2cf534fdc26cc3", "https://dummy.com:443",
		"b50a6e20d3733fb89631ae32385b3c85c533c560", []string{"0021"}, 21)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.SubmitWithCtx(context.Background(), Mainnet, stakeNode, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, stakeNode, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnstakeNodeWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unstakeNode, err := NewUnstakeNode("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.SubmitWithCtx(context.Background(), Mainnet, unstakeNode, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, unstakeNode, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}

func TestTransactionBuilder_SubmitUnjailNodeWithCtx(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	signer, err := signer.NewRandomSigner()
	c.NoError(err)

	txBuilder := NewTransactionBuilder(provider.NewProvider("https://dummy.com", []string{"https://dummy.com"}), signer)

	unjailNode, err := NewUnjailNode("b50a6e20d3733fb89631ae32385b3c85c533c560", "b50a6e20d3733fb89631ae32385b3c85c533c561")
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusOK, "../provider/samples/client_raw_tx.json")

	output, err := txBuilder.SubmitWithCtx(context.Background(), Mainnet, unjailNode, nil)
	c.NotEmpty(output)
	c.NoError(err)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", provider.ClientRawTXRoute),
		http.StatusInternalServerError, "../provider/samples/client_raw_tx.json")

	output, err = txBuilder.SubmitWithCtx(context.Background(), Mainnet, unjailNode, nil)
	c.Empty(output)
	c.Equal(provider.Err5xxOnConnection, err)
}
