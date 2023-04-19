// Package provider has the functions to do the RPC requests to a node in the protocol
package provider

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"

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

// Provider struct handler por JSON RPC provider
type Provider struct {
	rpcURL      string
	dispatchers []string
	client      *http.Client
}

// NewProvider returns Provider instance from input
func NewProvider(rpcURL string, dispatchers []string) *Provider {
	return &Provider{
		rpcURL:      rpcURL,
		dispatchers: dispatchers,
		client:      &http.Client{},
	}
}

func (p *Provider) getFinalRPCURL(rpcURL string, route V1RPCRoute) (string, error) {
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

func (p *Provider) doPostRequest(rpcURL string, params any, route V1RPCRoute) (*http.Response, error) {
	finalRPCURL, err := p.getFinalRPCURL(rpcURL, route)
	if err != nil {
		return nil, err
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", finalRPCURL, route), bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	output, err := p.client.Do(req)
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
	defer utils.CloseOrLog(body)

	if route == ClientRelayRoute {
		return errOnRelayRequest
	}

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
func (p *Provider) GetBalance(address string, options *GetBalanceOptions) (*big.Int, error) {
	params := map[string]any{
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
func (p *Provider) GetAccountTransactions(address string, options *GetAccountTransactionsOptions) (*GetAccountTransactionsOutput, error) {
	params := map[string]any{
		"address": address,
	}

	if options != nil {
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

// GetBlockTransactions returns transactions of given block
func (p *Provider) GetBlockTransactions(options *GetBlockTransactionsOptions) (*GetBlockTransactionsOutput, error) {
	params := map[string]any{}

	if options != nil {
		params["height"] = options.Height
		params["page"] = options.Page
		params["per_page"] = options.PerPage
		params["prove"] = options.Prove
		params["order"] = options.Order
	}

	rawOutput, err := p.doPostRequest("", params, QueryBlockTXsRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetBlockTransactionsOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// AddressType enum listing all address types
type AddressType string

const (
	// NodeType represents node type
	NodeType AddressType = "node"
	// AppType represents app type
	AppType AddressType = "app"
	// AccountType represents account type
	AccountType AddressType = "account"
)

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
func (p *Provider) GetType(address string, options *GetTypeOptions) (AddressType, error) {
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
func (p *Provider) SendTransaction(input *SendTransactionInput) (*SendTransactionOutput, error) {
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
func (p *Provider) GetBlock(blockNumber int) (*GetBlockOutput, error) {
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
func (p *Provider) GetTransaction(transactionHash string, options *GetTransactionOptions) (*GetTransactionOutput, error) {
	params := map[string]any{
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
func (p *Provider) GetBlockHeight() (int, error) {
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

// GetAllParams returns the params at the specified height
func (p *Provider) GetAllParams(options *GetAllParamsOptions) (*AllParams, error) {
	var height int
	if options != nil {
		height = options.Height
	}

	params := map[string]interface{}{
		"height": height,
	}

	rawOutput, err := p.doPostRequest("", params, QueryAllParamsRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	var allParams AllParams
	if err = json.Unmarshal(respBody, &allParams); err != nil {
		return nil, err
	}

	return &allParams, nil
}

// GetNodes returns a page of nodes known at the specified height and with options
// empty options returns all validators, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *Provider) GetNodes(options *GetNodesOptions) (*GetNodesOutput, error) {
	params := map[string]any{}

	if options != nil {
		params["height"] = options.Height
		params["opts"] = map[string]any{
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
func (p *Provider) GetNode(address string, options *GetNodeOptions) (*GetNodeOutput, error) {
	params := map[string]any{
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
func (p *Provider) GetApps(options *GetAppsOptions) (*GetAppsOutput, error) {
	params := map[string]any{}

	if options != nil {
		params["height"] = options.Height
		params["opts"] = map[string]any{
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
func (p *Provider) GetApp(address string, options *GetAppOptions) (*GetAppOutput, error) {
	params := map[string]any{
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
func (p *Provider) GetAccount(address string, options *GetAccountOptions) (*GetAccountOutput, error) {
	params := map[string]any{
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

// GetAccounts returns a page of accounts known at the specified height and with options
// empty options returns all accounts on last height, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *Provider) GetAccounts(options *GetAccountsOptions) (*GetAccountsOutput, error) {
	params := map[string]any{}

	if options != nil {
		params["height"] = options.Height
		params["page"] = options.Page
		params["per_page"] = options.PerPage
	}

	rawOutput, err := p.doPostRequest("", params, QueryAccountsRoute)

	defer closeOrLog(rawOutput)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(rawOutput.Body)
	if err != nil {
		return nil, err
	}

	output := GetAccountsOutput{}

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// Dispatch sends a dispatch request to the network and gets the nodes that will be servicing the requests for the session.
func (p *Provider) Dispatch(appPublicKey, chain string, options *DispatchRequestOptions) (*DispatchOutput, error) {
	if len(p.dispatchers) == 0 {
		return nil, ErrNoDispatchers
	}

	params := map[string]any{
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
func (p *Provider) Relay(rpcURL string, input *RelayInput, options *RelayRequestOptions) (*RelayOutput, error) {
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
