package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/pokt-foundation/pocket-go/internal/client"
	internalUtils "github.com/pokt-foundation/pocket-go/internal/utils"
)

var (
	// Err4xxOnConnection error when RPC responds with 4xx code
	Err4xxOnConnection = errors.New("rpc responded with 4xx")
	// Err5xxOnConnection error when RPC responds with 5xx code
	Err5xxOnConnection = errors.New("rpc responded with 5xx")
	// ErrUnexpectedResponse error when RPC responds with unexpected code
	ErrUnexpectedResponse = errors.New("rpc responded with unexpected code")
)

type JSONRPCProvider struct {
	rpcURL      string
	dispatchers []string
	client      *client.Client
}

func (p *JSONRPCProvider) getFinalRPCURL(rpcURL string, route V1RPCRoute) string {
	if rpcURL != "" {
		return rpcURL
	}

	if route == ClientDispatch {
		return p.dispatchers[int(math.Floor(rand.Float64()*100))%len(p.dispatchers)]
	}

	return p.rpcURL
}

func (p *JSONRPCProvider) doPostRequest(rpcURL string, params map[string]string, route V1RPCRoute) (*http.Response, error) {
	finalRPCURL := p.getFinalRPCURL(rpcURL, route)

	response, err := p.client.PostWithURLJSONParams(fmt.Sprintf("%s%s", finalRPCURL, route), params, http.Header{})
	if err != nil {
		return nil, err
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

	return nil, ErrUnexpectedResponse
}

func (p *JSONRPCProvider) GetBalance(address string) (*big.Int, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryBalance)
	if err != nil {
		return nil, err
	}

	defer internalUtils.CloseOrLog(rawResponse.Body)

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

func (p *JSONRPCProvider) GetTransactionCount(address string) (int, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryAccountTXs)
	if err != nil {
		return 0, err
	}

	defer internalUtils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return 0, err
	}

	response := queryAccountsTXsResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(response.TotalCount)
}

func returnType(appResponse queryAppResponse, nodeResponse queryNodeResponse) AddressType {
	if nodeResponse.ServiceURL == nil && appResponse.MaxRelays != nil {
		return App
	}

	if nodeResponse.ServiceURL != nil && appResponse.MaxRelays == nil {
		return Node
	}

	return Account
}

func (p *JSONRPCProvider) GetType(address string) (AddressType, error) {
	rawAppResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryApp)
	if err != nil {
		return "", err
	}

	defer internalUtils.CloseOrLog(rawAppResponse.Body)

	rawNodeResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryNode)
	if err != nil {
		return "", err
	}

	defer internalUtils.CloseOrLog(rawNodeResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawAppResponse.Body)
	if err != nil {
		return "", err
	}

	appResponse := queryAppResponse{}

	err = json.Unmarshal(bodyBytes, &appResponse)
	if err != nil {
		return "", err
	}

	bodyBytes, err = ioutil.ReadAll(rawNodeResponse.Body)
	if err != nil {
		return "", err
	}

	nodeResponse := queryNodeResponse{}

	err = json.Unmarshal(bodyBytes, &nodeResponse)
	if err != nil {
		return "", err
	}

	return returnType(appResponse, nodeResponse), nil
}

func (p *JSONRPCProvider) GetTransaction(address string) (*TransactionReponse, error) {
	rawResponse, err := p.doPostRequest("", map[string]string{
		"address": address,
	}, QueryTX)
	if err != nil {
		return nil, err
	}

	defer internalUtils.CloseOrLog(rawResponse.Body)

	bodyBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	response := TransactionReponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	if response.Transaction.Hash == "" {
		return nil, ErrUnexpectedResponse
	}

	return &response, nil
}
