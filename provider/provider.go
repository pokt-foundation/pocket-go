// Package provider has the functions to do the RPC requests to a node in the protocol
package provider

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pokt-foundation/pocket-go/utils"
	"github.com/pokt-foundation/utils-go/client"
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

	regexPatterns = []*regexp.Regexp{}

	errorStatusCodesMap = map[string]int{
		"context deadline exceeded (Client.Timeout exceeded while awaiting headers)": 408, // Request Timeout
		"connection reset by peer":            503, // Service Unavailable
		"no such host":                        404, // Not Found
		"network is unreachable":              503, // Service Unavailable
		"connection refused":                  502, // Bad Gateway
		"http: server closed idle connection": 499, // Client Closed Request (non-standard)
		"tls: handshake failure":              525, // SSL Handshake Failed (Cloudflare specific, non-standard)
		"i/o timeout":                         504, // Gateway Timeout
		"bad gateway":                         502, // Bad Gateway
		"service unavailable":                 503, // Service Unavailable
		"gateway timeout":                     504, // Gateway Timeout
	}
)

const (
	// DefaultStatusCode means the response was accepted but we don't know if it actually succeeded
	DefaultStatusCode = http.StatusAccepted
	// Use this contante to avoid the use of the hardcoded string result
	// result is the field present in a successful response
	resultText = "result"
)

// Provider struct handler por JSON RPC provider
type Provider struct {
	rpcURL      string
	dispatchers []string
	client      *client.Client
}

// NewProvider returns Provider instance from input
func NewProvider(rpcURL string, dispatchers []string) *Provider {
	return &Provider{
		rpcURL:      rpcURL,
		dispatchers: dispatchers,
		client:      client.NewDefaultClient(),
	}
}

func init() {
	regexPatterns = []*regexp.Regexp{
		regexp.MustCompile(`"code"\s*:\s*(\d+)`),       // Matches and captures any numeric status code after `"code":`
		regexp.MustCompile(`(\d+)\s+Not Found`),        // Matches and captures the status code from strings like `404 Not Found`
		regexp.MustCompile(`(\d+)\s+page not found`),   // Matches and captures the status code from strings like `404 page not found`
		regexp.MustCompile(`HTTP\/\d\.\d\s+(\d+)`),     // Matches and captures the status code from HTTP status lines like `HTTP/1.1 200`
		regexp.MustCompile(`"statusCode"\s*:\s*(\d+)`), // Matches and captures any numeric status code after `"statusCode":`
		regexp.MustCompile(`(\d+)\s+OK`),               // Matches and captures `200` in a response like `200 OK`
		regexp.MustCompile(`"statusCode"\s*:\s*(\d+)`), // Matches and captures any numeric status code after `"statusCode":`, added redundantly for clarity in different contexts
	}
}

// RequestConfigOpts are the optional values for request config
type RequestConfigOpts struct {
	Retries   int
	Timeout   time.Duration
	Transport http.RoundTripper
}

// UpdateRequestConfig updates the config used for RPC requests
func (p *Provider) UpdateRequestConfig(opts RequestConfigOpts) {
	p.client = client.NewCustomClientWithOptions(client.CustomClientOpts{
		Retries:   opts.Retries,
		Timeout:   opts.Timeout,
		Transport: opts.Transport,
	})
}

// ResetRequestConfigToDefault resets request config to default
func (p *Provider) ResetRequestConfigToDefault() {
	p.client = client.NewDefaultClient()
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

func (p *Provider) doPostRequest(ctx context.Context, rpcURL string, params any, route V1RPCRoute, headers http.Header) (*http.Response, error) {
	finalRPCURL, err := p.getFinalRPCURL(rpcURL, route)
	if err != nil {
		return nil, err
	}

	output, err := p.client.PostWithURLJSONParamsWithCtx(ctx, fmt.Sprintf("%s%s", finalRPCURL, route), params, headers)
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
func (p *Provider) GetBalance(address string, options *GetBalanceOptions) (*big.Int, error) {
	return p.GetBalanceWithCtx(context.Background(), address, options)
}

// GetBalanceWithCtx requests the balance of the specified address
func (p *Provider) GetBalanceWithCtx(ctx context.Context, address string, options *GetBalanceOptions) (*big.Int, error) {
	params := map[string]any{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryBalanceRoute, http.Header{})

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
	return p.GetAccountTransactionsWithCtx(context.Background(), address, options)
}

// GetAccountTransactionsWithCtx returns transactions of given address' account
func (p *Provider) GetAccountTransactionsWithCtx(ctx context.Context, address string, options *GetAccountTransactionsOptions) (*GetAccountTransactionsOutput, error) {
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

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryAccountTXsRoute, http.Header{})

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
	return p.GetBlockTransactionsWithCtx(context.Background(), options)
}

// GetBlockTransactionsWithCtx returns transactions of given block
func (p *Provider) GetBlockTransactionsWithCtx(ctx context.Context, options *GetBlockTransactionsOptions) (*GetBlockTransactionsOutput, error) {
	params := map[string]any{}

	if options != nil {
		params["height"] = options.Height
		params["page"] = options.Page
		params["per_page"] = options.PerPage
		params["prove"] = options.Prove
		params["order"] = options.Order
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryBlockTXsRoute, http.Header{})

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
	return p.GetTypeWithCtx(context.Background(), address, options)
}

// GetTypeWithCtx returns type of given address
func (p *Provider) GetTypeWithCtx(ctx context.Context, address string, options *GetTypeOptions) (AddressType, error) {
	var height int
	var errOutput *RPCError

	if options != nil {
		height = options.Height
	}

	_, appErr := p.GetAppWithCtx(ctx, address, &GetAppOptions{Height: height})
	if appErr != nil && !errors.As(appErr, &errOutput) {
		return "", appErr
	}

	_, nodeErr := p.GetNodeWithCtx(ctx, address, &GetNodeOptions{Height: height})
	if nodeErr != nil && !errors.As(nodeErr, &errOutput) {
		return "", nodeErr
	}

	return returnType(appErr, nodeErr), nil
}

// SendTransaction sends raw transaction to be relayed to a target address
func (p *Provider) SendTransaction(input *SendTransactionInput) (*SendTransactionOutput, error) {
	return p.SendTransactionWithCtx(context.Background(), input)
}

// SendTransactionWithCtx sends raw transaction to be relayed to a target address
func (p *Provider) SendTransactionWithCtx(ctx context.Context, input *SendTransactionInput) (*SendTransactionOutput, error) {
	rawOutput, err := p.doPostRequest(ctx, "", input, ClientRawTXRoute, http.Header{})

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
	return p.GetBlockWithCtx(context.Background(), blockNumber)
}

// GetBlockWithCtx returns the block structure at the specified height, height = 0 is used as latest
func (p *Provider) GetBlockWithCtx(ctx context.Context, blockNumber int) (*GetBlockOutput, error) {
	rawOutput, err := p.doPostRequest(ctx, "", map[string]int{
		"height": blockNumber,
	}, QueryBlockRoute, http.Header{})

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
	return p.GetTransactionWithCtx(context.Background(), transactionHash, options)
}

// GetTransactionWithCtx returns the transaction by the given transaction hash
func (p *Provider) GetTransactionWithCtx(ctx context.Context, transactionHash string, options *GetTransactionOptions) (*GetTransactionOutput, error) {
	params := map[string]any{
		"hash": transactionHash,
	}

	if options != nil {
		params["prove"] = options.Prove
	}

	rawOutput, err := p.doPostRequest(context.Background(), "", params, QueryTXRoute, http.Header{})

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

// GetNodeBlockHeight returns the blockheight reported by a specific node
func (p *Provider) GetNodeBlockHeight(ctx context.Context, url string) (int, error) {
	return p.getNodeBlockHeight(ctx, url)
}

// GetBlockHeight returns the current height
func (p *Provider) GetBlockHeight() (int, error) {
	return p.GetBlockHeightWithCtx(context.Background())
}

// GetBlockHeightWithCtx returns the current height
func (p *Provider) GetBlockHeightWithCtx(ctx context.Context) (int, error) {
	return p.getNodeBlockHeight(ctx, "")
}

func (p *Provider) getNodeBlockHeight(ctx context.Context, url string) (int, error) {
	rawOutput, err := p.doPostRequest(ctx, url, nil, QueryHeightRoute, http.Header{})
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
	return p.GetAllParamsWithCtx(context.Background(), options)
}

// GetAllParamsWithCtx returns the params at the specified height
func (p *Provider) GetAllParamsWithCtx(ctx context.Context, options *GetAllParamsOptions) (*AllParams, error) {
	var (
		height int
		url    string
	)

	if options != nil {
		height = options.Height
		url = options.URL
	}

	params := map[string]interface{}{
		"height": height,
	}

	rawOutput, err := p.doPostRequest(ctx, url, params, QueryAllParamsRoute, http.Header{})

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
	return p.GetNodesWithCtx(context.Background(), options)
}

// GetNodesWithCtx returns a page of nodes known at the specified height and with options
// empty options returns all validators, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *Provider) GetNodesWithCtx(ctx context.Context, options *GetNodesOptions) (*GetNodesOutput, error) {
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

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryNodesRoute, http.Header{})

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
	return p.GetNodeWithCtx(context.Background(), address, options)
}

// GetNodeWithCtx returns the node at the specified height, height = 0 is used as latest
func (p *Provider) GetNodeWithCtx(ctx context.Context, address string, options *GetNodeOptions) (*GetNodeOutput, error) {
	params := map[string]any{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryNodeRoute, http.Header{})

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
	return p.GetAppsWithCtx(context.Background(), options)
}

// GetAppsWithCtx returns a page of applications known at the specified height and staking status
// empty ("") staking_status returns all apps, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *Provider) GetAppsWithCtx(ctx context.Context, options *GetAppsOptions) (*GetAppsOutput, error) {
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

	headers := http.Header{}
	if options != nil && options.Headers != nil {
		headers = options.Headers
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryAppsRoute, headers)

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
	return p.GetAppWithCtx(context.Background(), address, options)
}

// GetAppWithCtx returns the app at the specified height, height = 0 is used as latest
func (p *Provider) GetAppWithCtx(ctx context.Context, address string, options *GetAppOptions) (*GetAppOutput, error) {
	params := map[string]any{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryAppRoute, http.Header{})

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
	return p.GetAccountWithCtx(context.Background(), address, options)
}

// GetAccountWithCtx returns account at the specified address
func (p *Provider) GetAccountWithCtx(ctx context.Context, address string, options *GetAccountOptions) (*GetAccountOutput, error) {
	params := map[string]any{
		"address": address,
	}

	if options != nil {
		params["height"] = options.Height
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryAccountRoute, http.Header{})

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
	return p.GetAccountsWithCtx(context.Background(), options)
}

// GetAccountsWithCtx returns a page of accounts known at the specified height and with options
// empty options returns all accounts on last height, page < 1 returns the first page, per_page < 1 returns 10000 elements per page
func (p *Provider) GetAccountsWithCtx(ctx context.Context, options *GetAccountsOptions) (*GetAccountsOutput, error) {
	params := map[string]any{}

	if options != nil {
		params["height"] = options.Height
		params["page"] = options.Page
		params["per_page"] = options.PerPage
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, QueryAccountsRoute, http.Header{})

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
	return p.DispatchWithCtx(context.Background(), appPublicKey, chain, options)
}

// DispatchWithCtx sends a dispatch request to the network and gets the nodes that will be servicing the requests for the session.
func (p *Provider) DispatchWithCtx(ctx context.Context, appPublicKey, chain string, options *DispatchRequestOptions) (*DispatchOutput, error) {
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

	headers := http.Header{}
	if options != nil && options.Headers != nil {
		headers = options.Headers
	}

	rawOutput, err := p.doPostRequest(ctx, "", params, ClientDispatchRoute, headers)

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
// Will always return with an output that includes the status code from the request
func (p *Provider) Relay(rpcURL string, input *RelayInput, options *RelayRequestOptions) (*RelayOutput, error) {
	return p.RelayWithCtx(context.Background(), rpcURL, input, options)
}

// RelayWithCtx does request to be relayed to a target blockchain
// Will always return with an output that includes the status code from the request
func (p *Provider) RelayWithCtx(ctx context.Context, rpcURL string, input *RelayInput, options *RelayRequestOptions) (*RelayOutput, error) {
	rawOutput, reqErr := p.doPostRequest(ctx, rpcURL, input, ClientRelayRoute, http.Header{})

	defer closeOrLog(rawOutput)

	statusCode := extractStatusFromRequest(rawOutput, reqErr)

	defaultOutput := &RelayOutput{
		StatusCode: statusCode,
	}

	if reqErr != nil && !errors.Is(reqErr, errOnRelayRequest) {
		return defaultOutput, reqErr
	}

	bodyBytes, err := io.ReadAll(rawOutput.Body)
	if err != nil {
		return defaultOutput, err
	}

	if errors.Is(reqErr, errOnRelayRequest) {
		return defaultOutput, parseRelayErrorOutput(bodyBytes, input.Proof.ServicerPubKey)
	}

	// The statusCode will be overwritten based on the response
	return parseRelaySuccesfulOutput(bodyBytes, statusCode)
}

func extractStatusFromRequest(rawOutput *http.Response, reqErr error) int {
	statusCode := DefaultStatusCode

	if reqErr != nil {
		for key, status := range errorStatusCodesMap {
			if strings.Contains(reqErr.Error(), key) { // This checks if the actual error contains the key string
				return status
			}
		}

		// If we got an error and we can't identify it as a known error, will be mark as if the server failed
		return http.StatusInternalServerError
	}

	if rawOutput.StatusCode != http.StatusOK {
		// If there's a response we'll use that as the status
		// NOTE: We know that nodes are manipulating the output, for this reason we'll ignore the status if it's ok
		statusCode = rawOutput.StatusCode
	}

	return statusCode
}

// TODO: Remove this function after the node responds back to us with a statusCode alongside with the response and the signature.
// Returns 202 if none of the pre-defined internal regexes matches any return values.
func extractStatusFromResponse(response string) int {
	for _, pattern := range regexPatterns {
		matches := pattern.FindStringSubmatch(response)
		if len(matches) > 1 {
			code, err := strconv.Atoi(matches[1])
			if err != nil || http.StatusText(code) == "" {
				continue
			}
			return code
		}
	}
	return DefaultStatusCode
}

func parseRelaySuccesfulOutput(bodyBytes []byte, requestStatusCode int) (*RelayOutput, error) {
	output := RelayOutput{
		StatusCode: requestStatusCode,
	}

	err := json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return &output, err
	}

	// Check if there's explicitly a result field, if there's on mark it as success, otherwise check what's the potential status.
	// for REST chain that doesn't return result in any of the call will be defaulted to 202 in extractStatusFromResponse
	if strings.Contains(output.Response, resultText) {
		output.StatusCode = http.StatusOK
	} else {
		output.StatusCode = extractStatusFromResponse(output.Response)
	}

	if !json.Valid([]byte(output.Response)) {
		return &output, ErrNonJSONResponse
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
