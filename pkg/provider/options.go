package provider

// GetBalanceOptions represents optional arguments for GetBalance request
type GetBalanceOptions struct {
	Height int
}

// GetNodeOptions represents optional arguments for GetNode request
type GetNodeOptions struct {
	Height int
}

// GetNodesOptions represents optional arguments for GetNodes request
type GetNodesOptions struct {
	StakingStatus StakingStatus
	Page          int
	PerPage       int
	BlockChain    string
	JailedStatus  JailedStatus
}

// GetAppOptions represents optional arguments for GetApp request
type GetAppOptions struct {
	Height int
}

// GetAccountOptions represents optional arguments for GetAccount request
type GetAccountOptions struct {
	Height int
}

// GetTypeOptions represents optional arguments for GetType request
type GetTypeOptions struct {
	Height int
}

// GetAppsOptions represents optional arguments for GetApps request
type GetAppsOptions struct {
	StakingStatus StakingStatus
	Page          int
	PerPage       int
	BlockChain    string
}

// DispatchRequestOptions represents optional arguments for Dispatch request
type DispatchRequestOptions struct {
	Height                       int
	RejectSelfSignedCertificates bool
}

// RelayRequestOptions represents optional arguments for Relay request
type RelayRequestOptions struct {
	RejectSelfSignedCertificates bool
}

// GetTransactionOptions represents the optional arguments for a GetTransaction request
type GetTransactionOptions struct {
	Prove bool
}
