package models

// RequestHash interface that represents a request hash
type RequestHash interface {
	GetPayload() RelayPayload
	GetMeta() RelayMeta
}
