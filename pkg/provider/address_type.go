package provider

// AddressType enum listing all address types
type AddressType string

const (
	// NodeType represents node type
	NodeType AddressType = "node"
	// AppType represents app type
	AppType AddressType = "app"
	// AccountType represents account type
	AccountType AddressType = "account"
)
