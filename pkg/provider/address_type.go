package provider

// AddressType enum listing all address types
type AddressType string

const (
	// Node represents node type
	Node AddressType = "node"
	// App represents app type
	App AddressType = "app"
	// Account represents account type
	Account AddressType = "account"
)
