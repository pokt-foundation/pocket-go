package models

// Node interface that represents a node
type Node interface {
	GetAddress() string
	GetChains() []string
	IsJailed() bool
	GetPublicKey() string
	GetServiceURL() string
	GetStakedTokens() string
	GetStatus() StakingStatus
	GetUnstakingTime() string
}
