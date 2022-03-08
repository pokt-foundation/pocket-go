package types

// App interface that represents an app
type App interface {
	GetAddress() string
	GetChains() []string
	IsJailed() bool
	GetMaxRelays() string
	GetPublicKey() string
	GetStakedTokens() string
	GetStatus() StackingStatus
}
