package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/pokt-foundation/pocket-go/pkg/client"
	"github.com/pokt-foundation/pocket-go/pkg/utils"
)

var (
	// Err4xxOnConnection error when RPC responds with 4xx code
	Err4xxOnConnection = errors.New("rpc responded with 4xx")
	// Err5xxOnConnection error when RPC responds with 5xx code
	Err5xxOnConnection = errors.New("rpc responded with 5xx")
	// ErrUnexpectedCodeOnConnection error when RPC responds with unexpected code
	ErrUnexpectedCodeOnConnection = errors.New("rpc responded with unexpected code")
	// ErrNoDispatchers error when dispatch call is requested with no dispatchers set
	ErrNoDispatchers = errors.New("no dispatchers")
	// ErrNonJSONResponse error when provider does not respond with a JSON
	ErrNonJSONResponse = errors.New("non JSON response")
)

// JSONRPCProvider struct handler por JSON RPC provider
type JSONRPCProvider struct {
	rpcURL      string
	dispatchers []string
	client      *client.Client
}

// NewJSONRPCProvider returns JSONRPCProvider instance from input
func NewJSONRPCProvider(rpcURL string, dispatchers []string, providerClient *client.Client) *JSONRPCProvider {
	return &JSONRPCProvider{
		rpcURL:      rpcURL,
		dispatchers: dispatchers,
		client:      providerClient,
	}
}

func (p *JSONRPCProvider) getFinalRPCURL(rpcURL string, route V1RPCRoute) string {
	if rpcURL != "" {
		return rpcURL
	}

	if route == ClientDispatchRoute {
		return p.dispatchers[rand.Intn(len(p.dispatchers))]
	}

	return p.rpcURL
}

func (p *JSONRPCProvider) doPostRequest(rpcURL string, params interface{}, route V1RPCRoute) (*http.Response, error) {
	finalRPCURL := p.getFinalRPCURL(rpcURL, route)

	response, err := p.client.PostWithURLJSONParams(fmt.Sprintf("%s%s", finalRPCURL, route), params, http.Header{})
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusBadRequest && route != ClientRelayRoute {
		return nil, returnRPCError(response.Body)
	}

	if string(response.Status[0]) == "4" {
		return response, Err4xxOnConnection
	}

	if string(response.Status[0]) == "5" {
		return response, Err5xxOnConnection
	}

	if string(response.Status[0]) == "2" {
		return response, nil
	}

	return nil, ErrUnexpectedCodeOnConnection
}

func returnRPCError(body io.ReadCloser) error {
	defer utils.CloseOrLog(body)

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	response := RPCError{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return err
	}

	return &response
}

// GetBalance requests the balance of the specified address
func (p *JSONRPCProvider) GetBalance(address string) (*big.Int, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryBalanceRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := queryBalanceResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return response.Balance, nil
}

func (p *JSONRPCProvider) queryAccountTXs(address string) (*queryAccountsTXsResponse, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryAccountTXsRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := queryAccountsTXsResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTransactionCount returns number of transactions sent by the given address
func (p *JSONRPCProvider) GetTransactionCount(address string) (int, error) {
	response, err := p.queryAccountTXs(address)
	if err != nil {
		return 0, err
	}

	return response.TotalCount, nil
}

func returnType(appErr, nodeErr error) AddressType {
	if nodeErr != nil && appErr == nil {
		return AppType
	}

	if nodeErr == nil && appErr != nil {
		return NodeType
	}

	return AccountType
}

// GetType returns type of given address
func (p *JSONRPCProvider) GetType(address string, options *GetTypeOptions) (AddressType, error) {
	var height int
	var errResponse *RPCError

	if options != nil {
		height = options.Height
	}

	_, appErr := p.GetApp(address, &GetAppOptions{Height: height})
	if appErr != nil && !errors.As(appErr, &errResponse) {
		return "", appErr
	}

	_, nodeErr := p.GetNode(address, &GetNodeOptions{Height: height})
	if nodeErr != nil && !errors.As(nodeErr, &errResponse) {
		return "", nodeErr
	}

	return returnType(appErr, nodeErr), nil
}

// SendTransaction sends raw transaction to be relayed to a target address
func (p *JSONRPCProvider) SendTransaction(signerAddress, signedTransaction string) (*SendTransactionResponse, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"address":       signerAddress,
		"raw_hex_bytes": signedTransaction,
	}, ClientRawTXRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := SendTransactionResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBlock returns the block structure at the specified height, height = 0 is used as latest
func (p *JSONRPCProvider) GetBlock(blockNumber int) (*GetBlockResponse, error) {
	rawResponse, err := p.doPostRequest("", map[string]int{
		"height": blockNumber,
	}, QueryBlockRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := GetBlockResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTransaction returns the transaction by the given transaction hash
func (p *JSONRPCProvider) GetTransaction(transactionHash string) (*GetTransactionResponse, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"hash": transactionHash,
	}, QueryTXRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := GetTransactionResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBlockNumber returns the current height
func (p *JSONRPCProvider) GetBlockNumber() (int, error) {
	rawResponse, err := p.doPostRequest("", nil, QueryHeightRoute)
	if err != nil {
		return 0, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return 0, err
	}

	response := queryHeightResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return 0, err
	}

	return response.Height, nil
}

// GetNodes returns a page of nodes known at the specified height and with options
// empty options returns all validators, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *JSONRPCProvider) GetNodes(height int, options *GetNodesOptions) (*GetNodesResponse, error) {
	params := map[string]interface{}{
		"height": height,
	}

	if options != nil {
		params["opts"] = map[string]interface{}{
			"staking_status": options.StakingStatus,
			"page":           options.Page,
			"per_page":       options.PerPage,
			"chain":          options.Chain,
			"jailed_status":  options.JailedStatus,
			"blockchain":     options.Blockchain,
		}
	}

	rawResponse, err := p.doPostRequest("", params, QueryNodesRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := GetNodesResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetNode returns the node at the specified height, height = 0 is used as latest
func (p *JSONRPCProvider) GetNode(address string, options *GetNodeOptions) (*GetNodeResponse, error) {
	params := map[string]interface{}{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawResponse, err := p.doPostRequest("", params, QueryNodeRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := GetNodeResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetApps returns a page of applications known at the specified height and staking status
// empty ("") staking_status returns all apps, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *JSONRPCProvider) GetApps(height int, options *GetAppsOptions) (*GetAppsResponse, error) {
	params := map[string]interface{}{
		"height": height,
	}

	if options != nil {
		params["opts"] = map[string]interface{}{
			"staking_status": options.StakingStatus,
			"page":           options.Page,
			"per_page":       options.PerPage,
			"chain":          options.Chain,
			"jailed_status":  options.JailedStatus,
			"blockchain":     options.Blockchain,
		}
	}

	rawResponse, err := p.doPostRequest("", params, QueryAppsRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := GetAppsResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetApp returns the app at the specified height, height = 0 is used as latest
func (p *JSONRPCProvider) GetApp(address string, options *GetAppOptions) (*GetAppResponse, error) {
	params := map[string]interface{}{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawResponse, err := p.doPostRequest("", params, QueryAppRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := GetAppResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetAccount returns account at the specified address
func (p *JSONRPCProvider) GetAccount(address string) (*GetAccountResponse, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryAccountRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := GetAccountResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetAccountWithTransactions returns account at the specified address with its performed transactions
func (p *JSONRPCProvider) GetAccountWithTransactions(address string) (*GetAccountWithTransactionsResponse, error) {
	accountResponse, err := p.GetAccount(address)
	if err != nil {
		return nil, err
	}

	transactionsResponse, err := p.queryAccountTXs(address)
	if err != nil {
		return nil, err
	}

	return &GetAccountWithTransactionsResponse{
		Account:      accountResponse,
		Transactions: transactionsResponse,
	}, nil
}

// Dispatch sends a dispatch request to the network and gets the nodes that will be servicing the requests for the session.
func (p *JSONRPCProvider) Dispatch(appPublicKey, chain string, sessionHeight int, options *DispatchRequestOptions) (*DispatchResponse, error) {
	if len(p.dispatchers) == 0 {
		return nil, ErrNoDispatchers
	}

	rawResponse, err := p.doPostRequest("", map[string]interface{}{
		"app_public_key": appPublicKey,
		"chain":          chain,
		"session_height": sessionHeight,
	}, ClientDispatchRoute)
	if err != nil {
		return nil, err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := DispatchResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Relay does request to be relayed to a target blockchain
func (p *JSONRPCProvider) Relay(rpcURL string, input *Relay, options *RelayRequestOptions) (string, error) {
	rawResponse, err := p.doPostRequest(rpcURL, input, ClientRelayRoute)
	if err != nil && rawResponse.StatusCode != http.StatusBadRequest {
		return "", err
	}

	defer utils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return "", err
	}

	if rawResponse.StatusCode == http.StatusBadRequest {
		return "", parseRelayErrorResponse(bodyBytes, input.Proof.ServicerPubKey)
	}

	if !json.Valid(bodyBytes) {
		return "", ErrNonJSONResponse
	}

	return string(bodyBytes), nil
}

func parseRelayErrorResponse(bodyBytes []byte, servicerPubKey string) error {
	response := RelayErrorResponse{}

	err := json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return err
	}

	return &RelayError{
		Code:           response.Error.Code,
		Codespace:      response.Error.Message,
		Message:        response.Error.Message,
		ServicerPubKey: servicerPubKey,
	}
}
