package types

// StackingStatus enum that represents staking status
type StackingStatus int

const (
	// Unstacked represents unstacked status
	Unstacked StackingStatus = iota
	// Unstacking represents unstacking status
	Unstacking
	// Stacked represents stacked status
	Stacked
)
