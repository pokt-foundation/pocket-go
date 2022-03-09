package models

// SessionHeader interface that represent session header
type SessionHeader interface {
	GetApplicationPubKey() string
	GetChain() string
	GetSessionBlockHeight() int
}
