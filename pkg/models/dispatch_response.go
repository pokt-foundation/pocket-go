package models

// DispatchResponse interface represents a dispatch response
type DispatchResponse interface {
	GetBlockHeight() int
	GetSession() Session
}
