package types

// DispatchResponse interface represents a dispatch response
type DispatchResponse interface {
	GetBlockHeight() int
	GetSession() Session
}
