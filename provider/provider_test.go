package provider

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/utils-go/mock-client"
	"github.com/stretchr/testify/require"
)

func TestRelayError_ErrorInterface(t *testing.T) {
	c := require.New(t)

	err := &RelayError{}

	var i any = err

	_, ok := i.(error)
	c.True(ok)
}

func TestProvider_GetBalance(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBalanceRoute), http.StatusOK, "samples/query_balance.json")

	balance, err := provider.GetBalance("pjog", &GetBalanceOptions{Height: 21})
	c.NoError(err)
	c.Equal(big.NewInt(1000000000), balance)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBalanceRoute), http.StatusBadRequest, "samples/error_response.json")

	balance, err = provider.GetBalance("pjog", nil)
	c.Equal("Request failed with code: 400 and message: dummy error", err.Error())
	c.Empty(balance)
}

func TestProvider_GetAccountTransactions(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountTXsRoute), http.StatusOK, "samples/query_account_txs.json")

	transactions, err := provider.GetAccountTransactions("pjog", &GetAccountTransactionsOptions{Prove: false})
	c.NoError(err)
	c.NotEmpty(transactions)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountTXsRoute), http.StatusInternalServerError, "samples/query_account_txs.json")

	transactions, err = provider.GetAccountTransactions("pjog", nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(transactions)
}

func TestProvider_GetBlockTransactions(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBlockTXsRoute), http.StatusOK, "samples/query_block_txs.json")

	transactions, err := provider.GetBlockTransactions(&GetBlockTransactionsOptions{Prove: false})
	c.NoError(err)
	c.NotEmpty(transactions)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBlockTXsRoute), http.StatusInternalServerError, "samples/query_block_txs.json")

	transactions, err = provider.GetBlockTransactions(nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(transactions)
}

func TestProvider_GetType(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusOK, "samples/query_app.json")
	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusOK, "samples/query_node.json")

	addressType, err := provider.GetType("pjog", &GetTypeOptions{Height: 21})
	c.NoError(err)
	c.Equal(AccountType, addressType)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusUnauthorized, "samples/query_node.json")

	addressType, err = provider.GetType("pjog", nil)
	c.Equal(Err4xxOnConnection, err)
	c.Empty(addressType)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusMultipleChoices, "samples/query_app.json")

	addressType, err = provider.GetType("pjog", &GetTypeOptions{Height: 21})
	c.Equal(ErrUnexpectedCodeOnConnection, err)
	c.Empty(addressType)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusOK, "samples/query_app.json")
	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusBadRequest, "samples/error_response.json")

	addressType, err = provider.GetType("pjog", nil)
	c.NoError(err)
	c.Equal(AppType, addressType)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusBadRequest, "samples/error_response.json")
	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusOK, "samples/query_node.json")

	addressType, err = provider.GetType("pjog", nil)
	c.NoError(err)
	c.Equal(NodeType, addressType)
}

func TestProvider_SendTransaction(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	transaction, err := provider.SendTransaction(&SendTransactionInput{Address: "pjog", RawHexBytes: "abcd"})
	c.Contains(err.Error(), "Post \"https://dummy.com/v1/client/rawtx\": no responder found")
	c.Empty(transaction)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRawTXRoute), http.StatusOK, "samples/client_raw_tx.json")

	transaction, err = provider.SendTransaction(&SendTransactionInput{Address: "pjog", RawHexBytes: "abcd"})
	c.NoError(err)
	c.NotEmpty(transaction)
}

func TestProvider_GetBlock(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBlockRoute), http.StatusOK, "samples/query_block.json")

	block, err := provider.GetBlock(21)
	c.NoError(err)
	c.NotEmpty(block)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBlockRoute), http.StatusInternalServerError, "samples/query_block.json")

	block, err = provider.GetBlock(21)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(block)
}

func TestProvider_GetTransaction(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryTXRoute), http.StatusOK, "samples/query_tx.json")

	transaction, err := provider.GetTransaction("abcd", &GetTransactionOptions{Prove: true})
	c.NoError(err)
	c.NotEmpty(transaction)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryTXRoute), http.StatusInternalServerError, "samples/query_tx.json")

	transaction, err = provider.GetTransaction("abcd", nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(transaction)
}

func TestProvider_GetBlockHeight(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryHeightRoute), http.StatusOK, "samples/query_height.json")

	blockNumber, err := provider.GetBlockHeight()
	c.NoError(err)
	c.Equal(21, blockNumber)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryHeightRoute), http.StatusInternalServerError, "samples/query_height.json")

	blockNumber, err = provider.GetBlockHeight()
	c.Equal(Err5xxOnConnection, err)
	c.Empty(blockNumber)
}

func TestProvider_GetAllParams(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAllParamsRoute), http.StatusOK, "samples/query_allparams.json")

	allParams, err := provider.GetAllParams(nil)
	c.NoError(err)

	relaysToTokensMultiplier, exists := allParams.NodeParams.Get("pos/RelaysToTokensMultiplier")
	c.True(exists)
	c.Equal("2109", relaysToTokensMultiplier)
}

func TestProvider_GetNodes(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodesRoute), http.StatusOK, "samples/query_nodes.json")

	nodes, err := provider.GetNodes(&GetNodesOptions{Page: 2})
	c.NoError(err)
	c.NotEmpty(nodes)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodesRoute), http.StatusInternalServerError, "samples/query_nodes.json")

	nodes, err = provider.GetNodes(nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(nodes)
}

func TestProvider_GetApps(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppsRoute), http.StatusOK, "samples/query_apps.json")

	apps, err := provider.GetApps(&GetAppsOptions{Page: 2})
	c.NoError(err)
	c.NotEmpty(apps)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppsRoute), http.StatusInternalServerError, "samples/query_apps.json")

	apps, err = provider.GetApps(nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(apps)
}

func TestProvider_GetNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusOK, "samples/query_node.json")

	node, err := provider.GetNode("pjog", &GetNodeOptions{Height: 2})
	c.NoError(err)
	c.NotEmpty(node)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusInternalServerError, "samples/query_node.json")

	node, err = provider.GetNode("pjog", nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(node)
}

func TestProvider_GetApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusOK, "samples/query_app.json")

	app, err := provider.GetApp("pjog", &GetAppOptions{Height: 2})
	c.NoError(err)
	c.NotEmpty(app)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusInternalServerError, "samples/query_app.json")

	app, err = provider.GetApp("pjog", nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(app)
}

func TestProvider_GetAccount(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountRoute), http.StatusOK, "samples/query_account.json")

	account, err := provider.GetAccount("pjog", &GetAccountOptions{Height: 21})
	c.NoError(err)
	c.NotEmpty(account)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountRoute), http.StatusInternalServerError, "samples/query_account.json")

	account, err = provider.GetAccount("pjog", nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(account)
}

func TestProvider_GetAccounts(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountsRoute), http.StatusOK, "samples/query_accounts.json")

	account, err := provider.GetAccounts(&GetAccountsOptions{Height: 21})
	c.NoError(err)
	c.NotEmpty(account)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountsRoute), http.StatusInternalServerError, "samples/query_accounts.json")

	account, err = provider.GetAccounts(nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(account)
}

func TestProvider_Dispatch(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := &Provider{
		rpcURL: "https://dummy.com",
	}

	dispatch, err := provider.Dispatch("pjog", "abcd", &DispatchRequestOptions{Height: 21})
	c.Equal(ErrNoDispatchers, err)
	c.Empty(dispatch)

	provider.dispatchers = []string{"https://dummy.com"}

	provider.UpdateRequestConfig(0, 5*time.Second)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientDispatchRoute), http.StatusOK, "samples/client_dispatch.json")

	dispatch, err = provider.Dispatch("pjog", "abcd", nil)
	c.NoError(err)
	c.NotEmpty(dispatch)

	provider.ResetRequestConfigToDefault()

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientDispatchRoute), http.StatusInternalServerError, "samples/client_dispatch.json")

	dispatch, err = provider.Dispatch("pjog", "abcd", &DispatchRequestOptions{Height: 21})
	c.Equal(Err5xxOnConnection, err)
	c.Empty(dispatch)
}

func TestProvider_Relay(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewProvider("https://dummy.com", []string{"https://dummy.com"})

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRelayRoute), http.StatusOK, "samples/client_relay.json")

	relay, err := provider.Relay("https://dummy.com", &RelayInput{}, nil)
	c.NoError(err)
	c.NotEmpty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRelayRoute), http.StatusInternalServerError, "samples/client_relay.json")

	relay, err = provider.Relay("https://dummy.com", &RelayInput{}, nil)
	c.Equal(Err5xxOnConnection, err)
	c.False(IsErrorCode(EmptyPayloadDataError, err))
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRelayRoute), http.StatusBadRequest, "samples/client_relay_error.json")

	relay, err = provider.Relay("https://dummy.com", &RelayInput{Proof: &RelayProof{ServicerPubKey: "PJOG"}}, nil)
	c.Equal("Request failed with code: 25, codespace: pocketcore and message: the payload data of the relay request is empty\nWith ServicerPubKey: PJOG", err.Error())
	c.True(IsErrorCode(EmptyPayloadDataError, err))
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRelayRoute), http.StatusOK, "samples/client_relay_non_json.json")

	relay, err = provider.Relay("https://dummy.com", &RelayInput{}, nil)
	c.Equal(ErrNonJSONResponse, err)
	c.Empty(relay)
}
