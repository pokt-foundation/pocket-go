package types

// Account interface that represents an Account
type Account interface {
	GetAddress() string
	GetBalance() string
	GetPublicKey() string
}
