package provider

import "time"

// GetBalanceOptions represents optional arguments for GetBalance request
type GetBalanceOptions struct {
	Height         int
	RequestOptions *RequestOptions
}

func (o *GetBalanceOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetNodeOptions represents optional arguments for GetNode request
type GetNodeOptions struct {
	Height         int
	RequestOptions *RequestOptions
}

func (o *GetNodeOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetNodesOptions represents optional arguments for GetNodes request
type GetNodesOptions struct {
	StakingStatus  StakingStatus
	Page           int
	PerPage        int
	BlockChain     string
	JailedStatus   JailedStatus
	RequestOptions *RequestOptions
}

func (o *GetNodesOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetAccountTransactionsOptions represents optional arguments for GetAccountTransactions request
type GetAccountTransactionsOptions struct {
	Height         int
	Page           int
	PerPage        int
	Prove          bool
	Received       bool
	Order          Order
	RequestOptions *RequestOptions
}

func (o *GetAccountTransactionsOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetTransactionCountOptions represents optional arguments for GetTransactionCount request
type GetTransactionCountOptions struct {
	Height         int
	Received       bool
	RequestOptions *RequestOptions
}

// GetAppOptions represents optional arguments for GetApp request
type GetAppOptions struct {
	Height         int
	RequestOptions *RequestOptions
}

func (o *GetAppOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetAccountOptions represents optional arguments for GetAccount request
type GetAccountOptions struct {
	Height         int
	RequestOptions *RequestOptions
}

func (o *GetAccountOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetTypeOptions represents optional arguments for GetType request
type GetTypeOptions struct {
	Height         int
	RequestOptions *RequestOptions
}

// GetAppsOptions represents optional arguments for GetApps request
type GetAppsOptions struct {
	StakingStatus  StakingStatus
	Page           int
	PerPage        int
	BlockChain     string
	RequestOptions *RequestOptions
}

func (o *GetAppsOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// DispatchRequestOptions represents optional arguments for Dispatch request
type DispatchRequestOptions struct {
	Height                       int
	RejectSelfSignedCertificates bool
	RequestOptions               *RequestOptions
}

func (o *DispatchRequestOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// RelayRequestOptions represents optional arguments for Relay request
type RelayRequestOptions struct {
	RejectSelfSignedCertificates bool
	RequestOptions               *RequestOptions
}

func (o *RelayRequestOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetTransactionOptions represents the optional arguments for a GetTransaction request
type GetTransactionOptions struct {
	Prove          bool
	RequestOptions *RequestOptions
}

func (o *GetTransactionOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// SendTransactionOptions represents optional arguments for a SendTransaction request
type SendTransactionOptions struct {
	RequestOptions *RequestOptions
}

func (o *SendTransactionOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetBlockOptions represents optional arguments for a GetBlock request
type GetBlockOptions struct {
	RequestOptions *RequestOptions
}

func (o *GetBlockOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// GetBlockHeightOptions represents optional arguments for a GetBlockHeight request
type GetBlockHeightOptions struct {
	RequestOptions *RequestOptions
}

func (o *GetBlockHeightOptions) getRequestOptions() *RequestOptions {
	if o == nil {
		return nil
	}

	return o.RequestOptions
}

// RequestOptions represents optional arguments related to the HTTP request itself
type RequestOptions struct {
	HTTPTimeout time.Duration
	HTTPRetries int
}

type requester interface {
	getRequestOptions() *RequestOptions
}
