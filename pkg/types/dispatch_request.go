package types

// DispatchRequest interface represents a dispatch request
type DispatchRequest interface {
	GetSessionHeader() SessionHeader
}
