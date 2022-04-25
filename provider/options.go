package provider

// JailedStatus enum that represents jailed status
type JailedStatus int

const (
	// Jailed status is when a node has been jailed due to missing a determined amount of blocks and/or byzantine behavior and thus cannot serve relays nor participate in consensus
	Jailed JailedStatus = iota + 1
	// Unjailed status is when a node is not jailed and thus can serve relays
	Unjailed
)

// StakingStatus enum that represents staking status
type StakingStatus int

const (
	// Unstaked represents unstaked status
	Unstaked StakingStatus = iota
	// Unstaking represents unstaking status
	Unstaking
	// Staked represents staked status
	Staked
)

// Order enum that represents the order which RPC requests should return their outputs
type Order string

const (
	// DescendantOrder represents greater to lower order
	DescendantOrder Order = "desc"
	// AscendantOrder represents lower to greater order
	AscendantOrder Order = "asc"
)

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

// GetAccountTransactionsOptions represents optional arguments for GetAccountTransactions request
type GetAccountTransactionsOptions struct {
	Height   int
	Page     int
	PerPage  int
	Prove    bool
	Received bool
	Order    Order
}

// GetTransactionCountOptions represents optional arguments for GetTransactionCount request
type GetTransactionCountOptions struct {
	Height   int
	Received bool
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
