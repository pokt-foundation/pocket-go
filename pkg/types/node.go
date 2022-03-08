package types

// Node interface that represents a node
type Node interface {
	GetAddress() string
	GetChains() []string
	IsJailed() bool
	GetPublicKey() string
	GetServiceURL() string
	GetStakedTokens() string
	GetStatus() StackingStatus
	GetUnstakingTime() string
}
