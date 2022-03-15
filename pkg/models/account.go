package models

// Account interface that represents a Pocket Network Protocol Account
type Account interface {
	GetAddress() string
	GetBalance() string
	GetPublicKey() string
	GetPrivateKey() string
}
