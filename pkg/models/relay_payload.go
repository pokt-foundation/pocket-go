package models

// RelayPayload interface that represents
type RelayPayload interface {
	GetData() string
	GetMethod() string
	GetPath() string
	GetHeaders() RelayHeaders
}
