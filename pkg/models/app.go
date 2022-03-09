package types

// App interface that represents an app
type App interface {
	GetAddress() string
	GetChains() []string
	GetMaxRelays() string
	GetPublicKey() string
	GetStakedTokens() string
	GetStatus() StakingStatus
}
