package provider

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
