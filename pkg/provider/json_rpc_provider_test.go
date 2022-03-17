package provider

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/pocket-go/internal/client"
	"github.com/pokt-foundation/pocket-go/internal/mock-client"
	"github.com/stretchr/testify/require"
)

func TestJSONRPCProvider_ProviderInterface(t *testing.T) {
	c := require.New(t)

	provider := &JSONRPCProvider{}

	var i interface{} = provider

	_, ok := i.(Provider)
	c.True(ok)
}

func TestJSONRPCProvider_GetBalance(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBalanceRoute), http.StatusOK, "samples/query_balance.json")

	balance, err := provider.GetBalance("pjog")
	c.NoError(err)
	c.Equal(big.NewInt(1000000000), balance)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBalanceRoute), http.StatusBadRequest, "samples/query_balance.json")

	balance, err = provider.GetBalance("pjog")
	c.Equal(Err4xxOnConnection, err)
	c.Empty(balance)
}

func TestJSONRPCProvider_GetTransactionCount(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountTXsRoute), http.StatusOK, "samples/query_account_txs.json")

	count, err := provider.GetTransactionCount("pjog")
	c.NoError(err)
	c.Equal(21, count)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountTXsRoute), http.StatusInternalServerError, "samples/query_account_txs.json")

	count, err = provider.GetTransactionCount("pjog")
	c.Equal(Err5xxOnConnection, err)
	c.Empty(count)
}

func TestJSONRPCProvider_GetType(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusOK, "samples/query_app.json")
	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusOK, "samples/query_node.json")

	addressType, err := provider.GetType("pjog")
	c.NoError(err)
	c.Equal(AccountType, addressType)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusBadRequest, "samples/query_node.json")

	addressType, err = provider.GetType("pjog")
	c.Equal(Err4xxOnConnection, err)
	c.Empty(addressType)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusMultipleChoices, "samples/query_app.json")

	addressType, err = provider.GetType("pjog")
	c.Equal(ErrUnexpectedCodeOnConnection, err)
	c.Empty(addressType)
}

func Test_returnType(t *testing.T) {
	c := require.New(t)

	maxRelays := 12
	app := &GetAppResponse{
		MaxRelays: &maxRelays,
	}

	node := &GetNodeResponse{}

	c.Equal(AppType, returnType(app, node))

	app.MaxRelays = nil

	serviceURL := "https://dummy.com"
	node.ServiceURL = &serviceURL

	c.Equal(NodeType, returnType(app, node))
}

func TestJSONRPCProvider_SendTransaction(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	transaction, err := provider.SendTransaction("pjog", "abcd")
	c.Contains(err.Error(), "Post \"https://dummy.com/v1/client/rawtx\": no responder found")
	c.Empty(transaction)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRawTXRoute), http.StatusOK, "samples/client_raw_tx.json")

	transaction, err = provider.SendTransaction("pjog", "abcd")
	c.NoError(err)
	c.NotEmpty(transaction)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRawTXRoute), http.StatusOK, "samples/client_raw_tx_unexpected.json")

	transaction, err = provider.SendTransaction("pjog", "abcd")
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(transaction)
}

func TestJSONRPCProvider_GetBlock(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBlockRoute), http.StatusOK, "samples/query_block.json")

	block, err := provider.GetBlock(21)
	c.NoError(err)
	c.NotEmpty(block)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBlockRoute), http.StatusInternalServerError, "samples/query_block.json")

	block, err = provider.GetBlock(21)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(block)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryBlockRoute), http.StatusOK, "samples/query_block_unexpected.json")

	block, err = provider.GetBlock(21)
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(block)
}

func TestJSONRPCProvider_GetTransaction(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryTXRoute), http.StatusOK, "samples/query_tx.json")

	transaction, err := provider.GetTransaction("abcd")
	c.NoError(err)
	c.NotEmpty(transaction)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryTXRoute), http.StatusInternalServerError, "samples/query_tx.json")

	transaction, err = provider.GetTransaction("abcd")
	c.Equal(Err5xxOnConnection, err)
	c.Empty(transaction)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryTXRoute), http.StatusOK, "samples/query_tx_unexpected.json")

	transaction, err = provider.GetTransaction("abcd")
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(transaction)
}

func TestJSONRPCProvider_GetBlockNumber(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryHeightRoute), http.StatusOK, "samples/query_height.json")

	blockNumber, err := provider.GetBlockNumber()
	c.NoError(err)
	c.Equal(21, blockNumber)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryHeightRoute), http.StatusInternalServerError, "samples/query_height.json")

	blockNumber, err = provider.GetBlockNumber()
	c.Equal(Err5xxOnConnection, err)
	c.Empty(blockNumber)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryHeightRoute), http.StatusOK, "samples/query_height_unexpected.json")

	blockNumber, err = provider.GetBlockNumber()
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(blockNumber)
}

func TestJSONRPCProvider_GetNodes(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodesRoute), http.StatusOK, "samples/query_nodes.json")

	nodes, err := provider.GetNodes(21, &GetNodesOptions{Page: 2})
	c.NoError(err)
	c.NotEmpty(nodes)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodesRoute), http.StatusInternalServerError, "samples/query_nodes.json")

	nodes, err = provider.GetNodes(21, &GetNodesOptions{Page: 2})
	c.Equal(Err5xxOnConnection, err)
	c.Empty(nodes)
}

func TestJSONRPCProvider_GetApps(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppsRoute), http.StatusOK, "samples/query_apps.json")

	apps, err := provider.GetApps(21, &GetAppsOptions{Page: 2})
	c.NoError(err)
	c.NotEmpty(apps)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppsRoute), http.StatusInternalServerError, "samples/query_apps.json")

	apps, err = provider.GetApps(21, &GetAppsOptions{Page: 2})
	c.Equal(Err5xxOnConnection, err)
	c.Empty(apps)
}

func TestJSONRPCProvider_GetNode(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusOK, "samples/query_node.json")

	node, err := provider.GetNode("pjog", &GetNodeOptions{Height: 2})
	c.NoError(err)
	c.NotEmpty(node)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusInternalServerError, "samples/query_node.json")

	node, err = provider.GetNode("pjog", &GetNodeOptions{Height: 2})
	c.Equal(Err5xxOnConnection, err)
	c.Empty(node)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryNodeRoute), http.StatusOK, "samples/query_node_unexpected.json")

	node, err = provider.GetNode("pjog", &GetNodeOptions{Height: 2})
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(node)
}

func TestJSONRPCProvider_GetApp(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusOK, "samples/query_app.json")

	app, err := provider.GetApp("pjog", &GetAppOptions{Height: 2})
	c.NoError(err)
	c.NotEmpty(app)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusInternalServerError, "samples/query_app.json")

	app, err = provider.GetApp("pjog", &GetAppOptions{Height: 2})
	c.Equal(Err5xxOnConnection, err)
	c.Empty(app)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAppRoute), http.StatusOK, "samples/query_app_unexpected.json")

	app, err = provider.GetApp("pjog", &GetAppOptions{Height: 2})
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(app)
}

func TestJSONRPCProvider_GetAccount(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountRoute), http.StatusOK, "samples/query_account.json")

	account, err := provider.GetAccount("pjog")
	c.NoError(err)
	c.NotEmpty(account)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountRoute), http.StatusInternalServerError, "samples/query_account.json")

	account, err = provider.GetAccount("pjog")
	c.Equal(Err5xxOnConnection, err)
	c.Empty(account)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountRoute), http.StatusOK, "samples/query_account_unexpected.json")

	account, err = provider.GetAccount("pjog")
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(account)
}

func TestJSONRPCProvider_GetAccountWithTransactions(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountRoute), http.StatusOK, "samples/query_account.json")
	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountTXsRoute), http.StatusOK, "samples/query_account_txs.json")

	account, err := provider.GetAccountWithTransactions("pjog")
	c.NoError(err)
	c.NotEmpty(account)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountTXsRoute), http.StatusOK, "samples/query_account_txs_unexpected.json")

	account, err = provider.GetAccountWithTransactions("pjog")
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(account)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", QueryAccountRoute), http.StatusInternalServerError, "samples/query_account.json")

	account, err = provider.GetAccountWithTransactions("pjog")
	c.Equal(Err5xxOnConnection, err)
	c.Empty(account)
}

func TestJSONRPCProvider_Dispatch(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	providerClient := client.NewDefaultClient()

	provider := &JSONRPCProvider{
		rpcURL: "https://dummy.com",
		client: providerClient,
	}

	dispatch, err := provider.Dispatch("pjog", "abcd", 21, nil)
	c.Equal(ErrNoDispatchers, err)
	c.Empty(dispatch)

	provider.dispatchers = []string{"https://dummy.com"}

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientDispatchRoute), http.StatusOK, "samples/client_dispatch.json")

	dispatch, err = provider.Dispatch("pjog", "abcd", 21, nil)
	c.NoError(err)
	c.NotEmpty(dispatch)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientDispatchRoute), http.StatusInternalServerError, "samples/client_dispatch.json")

	dispatch, err = provider.Dispatch("pjog", "abcd", 21, nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(dispatch)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientDispatchRoute), http.StatusOK, "samples/client_dispatch_unexpected.json")

	dispatch, err = provider.Dispatch("pjog", "abcd", 21, nil)
	c.Equal(ErrUnexpectedResponse, err)
	c.Empty(dispatch)
}

func TestJSONRPCProvider_Relay(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := NewJSONRPCProvider("https://dummy.com", []string{"https://dummy.com"}, client.NewDefaultClient())

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRelayRoute), http.StatusOK, "samples/client_relay.json")

	relay, err := provider.Relay("https://dummy.com", &Relay{}, nil)
	c.NoError(err)
	c.NotEmpty(relay.SuccessfulResponse)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRelayRoute), http.StatusInternalServerError, "samples/client_relay.json")

	relay, err = provider.Relay("https://dummy.com", &Relay{}, nil)
	c.Equal(Err5xxOnConnection, err)
	c.Empty(relay)

	mock.AddMockedResponseFromFile(http.MethodPost, fmt.Sprintf("%s%s", "https://dummy.com", ClientRelayRoute), http.StatusBadRequest, "samples/client_relay_error.json")

	relay, err = provider.Relay("https://dummy.com", &Relay{}, nil)
	c.NoError(err)
	c.NotEmpty(relay.ErrorResponse)
}
