package types

// Session interface that represents a session
type Session interface {
	GetBlockHeight() int
	GetHeader() SessionHeader
	GetKey() string
	GetNodes() []Node
}
