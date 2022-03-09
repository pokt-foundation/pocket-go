package models

// PocketAAT interface that represents a Pocket Application Authentication Token
type PocketAAT interface {
	GetVersion() string
	GetClientPublicKey() string
	GetApplicationPublicKey() string
	GetApplicationSignature() string
}
