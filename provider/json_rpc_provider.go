package provider

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/pokt-foundation/pocket-go/internal/client"
	"github.com/pokt-foundation/pocket-go/utils"
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

	errOnRelayRequest = errors.New("error on relay request")
)

// JSONRPCProvider struct handler por JSON RPC provider
type JSONRPCProvider struct {
	rpcURL      string
	dispatchers []string
	client      *client.Client
}

// NewJSONRPCProvider returns JSONRPCProvider instance from input
func NewJSONRPCProvider(rpcURL string, dispatchers []string) *JSONRPCProvider {
	return &JSONRPCProvider{
		rpcURL:      rpcURL,
		dispatchers: dispatchers,
		client:      client.NewDefaultClient(),
	}
}

// UpdateRequestConfig updates retries and timeout used for RPC requests
func (p *JSONRPCProvider) UpdateRequestConfig(retries int, timeout time.Duration) {
	p.client = client.NewCustomClient(retries, timeout)
}

// ResetRequestConfigToDefault resets request config to default
func (p *JSONRPCProvider) ResetRequestConfigToDefault() {
	p.client = client.NewDefaultClient()
}

func (p *JSONRPCProvider) getFinalRPCURL(rpcURL string, route V1RPCRoute) (string, error) {
	if rpcURL != "" {
		return rpcURL, nil
	}

	if route == ClientDispatchRoute {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(p.dispatchers))))
		if err != nil {
			return "", err
		}

		return p.dispatchers[index.Int64()], nil
	}

	return p.rpcURL, nil
}

func (p *JSONRPCProvider) doPostRequest(rpcURL string, params interface{}, route V1RPCRoute) (*http.Response, error) {
	finalRPCURL, err := p.getFinalRPCURL(rpcURL, route)
	if err != nil {
		return nil, err
	}

	output, err := p.client.PostWithURLJSONParams(fmt.Sprintf("%s%s", finalRPCURL, route), params, http.Header{})
	if err != nil {
		return nil, err
	}

	if output.StatusCode == http.StatusBadRequest {
		return output, returnRPCError(route, output.Body)
	}

	if string(output.Status[0]) == "4" {
		return output, Err4xxOnConnection
	}

	if string(output.Status[0]) == "5" {
		return output, Err5xxOnConnection
	}

	if string(output.Status[0]) == "2" {
		return output, nil
	}

	return nil, ErrUnexpectedCodeOnConnection
}

func returnRPCError(route V1RPCRoute, body io.ReadCloser) error {
	if route == ClientRelayRoute {
		return errOnRelayRequest
	}

	defer utils.CloseOrLog(body)

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	output := RPCError{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return err
	}

	return &output
}

// GetBalance requests the balance of the specified address
func (p *JSONRPCProvider) GetBalance(address string, options *GetBalanceOptions) (*big.Int, error) {
	params := map[string]interface{}{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest("", params, QueryBalanceRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := queryBalanceOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return output.Balance, nil
}

// GetAccountTransactions returns transactions of given address' account
func (p *JSONRPCProvider) GetAccountTransactions(address string, options *GetAccountTransactionsOptions) (*GetAccountTransactionsOutput, error) {
	params := map[string]interface{}{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
		params["page"] = options.Page
		params["per_page"] = options.PerPage
		params["prove"] = options.Prove
		params["received"] = options.Received
		params["order"] = options.Order
	}

	rawOutput, err := p.doPostRequest("", params, QueryAccountTXsRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetAccountTransactionsOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetTransactionCount returns number of transactions sent by the given address
func (p *JSONRPCProvider) GetTransactionCount(address string, options *GetTransactionCountOptions) (int, error) {
	optionsToSend := &GetAccountTransactionsOptions{}
	currentPage := 1
	totalCount := 0

	if options != nil {
		optionsToSend.Height = options.Height
		optionsToSend.Received = options.Received
	}

	for {
		optionsToSend.Page = currentPage

		output, err := p.GetAccountTransactions(address, optionsToSend)
		if err != nil {
			return 0, err
		}

		if output.TotalCount == 0 {
			break
		}

		totalCount += output.TotalCount
		currentPage++
	}

	return totalCount, nil
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
	var errOutput *RPCError

	if options != nil {
		height = options.Height
	}

	_, appErr := p.GetApp(address, &GetAppOptions{Height: height})
	if appErr != nil && !errors.As(appErr, &errOutput) {
		return "", appErr
	}

	_, nodeErr := p.GetNode(address, &GetNodeOptions{Height: height})
	if nodeErr != nil && !errors.As(nodeErr, &errOutput) {
		return "", nodeErr
	}

	return returnType(appErr, nodeErr), nil
}

// SendTransaction sends raw transaction to be relayed to a target address
func (p *JSONRPCProvider) SendTransaction(input *SendTransactionInput) (*SendTransactionOutput, error) {
	rawOutput, err := p.doPostRequest("", input, ClientRawTXRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := SendTransactionOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetBlock returns the block structure at the specified height, height = 0 is used as latest
func (p *JSONRPCProvider) GetBlock(blockNumber int) (*GetBlockOutput, error) {
	rawOutput, err := p.doPostRequest("", map[string]int{
		"height": blockNumber,
	}, QueryBlockRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetBlockOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetTransaction returns the transaction by the given transaction hash
func (p *JSONRPCProvider) GetTransaction(transactionHash string, options *GetTransactionOptions) (*GetTransactionOutput, error) {
	params := map[string]interface{}{
		"hash": transactionHash,
	}

	if options != nil {
		params["prove"] = options.Prove
	}

	rawOutput, err := p.doPostRequest("", params, QueryTXRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetTransactionOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetBlockHeight returns the current height
func (p *JSONRPCProvider) GetBlockHeight() (int, error) {
	rawOutput, err := p.doPostRequest("", nil, QueryHeightRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return 0, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return 0, err
	}

	output := queryHeightOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return 0, err
	}

	return output.Height, nil
}

// GetNodes returns a page of nodes known at the specified height and with options
// empty options returns all validators, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *JSONRPCProvider) GetNodes(height int, options *GetNodesOptions) (*GetNodesOutput, error) {
	params := map[string]interface{}{
		"height": height,
	}

	if options != nil {
		params["opts"] = map[string]interface{}{
			"staking_status": options.StakingStatus,
			"page":           options.Page,
			"per_page":       options.PerPage,
			"blockchain":     options.BlockChain,
			"jailed_status":  options.JailedStatus,
		}
	}

	rawOutput, err := p.doPostRequest("", params, QueryNodesRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetNodesOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetNode returns the node at the specified height, height = 0 is used as latest
func (p *JSONRPCProvider) GetNode(address string, options *GetNodeOptions) (*GetNodeOutput, error) {
	params := map[string]interface{}{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest("", params, QueryNodeRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetNodeOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetApps returns a page of applications known at the specified height and staking status
// empty ("") staking_status returns all apps, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *JSONRPCProvider) GetApps(height int, options *GetAppsOptions) (*GetAppsOutput, error) {
	params := map[string]interface{}{
		"height": height,
	}

	if options != nil {
		params["opts"] = map[string]interface{}{
			"staking_status": options.StakingStatus,
			"page":           options.Page,
			"per_page":       options.PerPage,
			"blockchain":     options.BlockChain,
		}
	}

	rawOutput, err := p.doPostRequest("", params, QueryAppsRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetAppsOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetApp returns the app at the specified height, height = 0 is used as latest
func (p *JSONRPCProvider) GetApp(address string, options *GetAppOptions) (*GetAppOutput, error) {
	params := map[string]interface{}{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest("", params, QueryAppRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetAppOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// GetAccount returns account at the specified address
func (p *JSONRPCProvider) GetAccount(address string, options *GetAccountOptions) (*GetAccountOutput, error) {
	params := map[string]interface{}{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest("", params, QueryAccountRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetAccountOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// Dispatch sends a dispatch request to the network and gets the nodes that will be servicing the requests for the session.
func (p *JSONRPCProvider) Dispatch(appPublicKey, chain string, options *DispatchRequestOptions) (*DispatchOutput, error) {
	if len(p.dispatchers) == 0 {
		return nil, ErrNoDispatchers
	}

	params := map[string]interface{}{
		"app_public_key": appPublicKey,
		"chain":          chain,
	}

	if options != nil {
		params["session_height"] = options.Height
	}

	rawOutput, err := p.doPostRequest("", params, ClientDispatchRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := DispatchOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// Relay does request to be relayed to a target blockchain
func (p *JSONRPCProvider) Relay(rpcURL string, input *Relay, options *RelayRequestOptions) (*RelayOutput, error) {
	rawOutput, reqErr := p.doPostRequest(rpcURL, input, ClientRelayRoute)

	defer closeOrLog(rawOutput)

	if reqErr != nil && !errors.Is(reqErr, errOnRelayRequest) {
		return nil, reqErr
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	if errors.Is(reqErr, errOnRelayRequest) {
		return nil, parseRelayErrorOutput(bodyBytes, input.Proof.ServicerPubKey)
	}

	return parseRelaySuccesfulOutput(bodyBytes)
}

func parseRelaySuccesfulOutput(bodyBytes []byte) (*RelayOutput, error) {
	output := RelayOutput{}

	err := json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	if !json.Valid([]byte(output.Response)) {
		return nil, ErrNonJSONResponse
	}

	return &output, nil
}

func parseRelayErrorOutput(bodyBytes []byte, servicerPubKey string) error {
	output := RelayErrorOutput{}

	err := json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return err
	}

	return &RelayError{
		Code:           output.Error.Code,
		Codespace:      output.Error.Codespace,
		Message:        output.Error.Message,
		ServicerPubKey: servicerPubKey,
	}
}

func closeOrLog(response *http.Response) {
	if response != nil {
		io.Copy(ioutil.Discard, response.Body)
		utils.CloseOrLog(response.Body)
	}
}
