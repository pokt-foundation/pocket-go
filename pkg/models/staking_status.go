package models

// StakingStatus enum that represents staking status
type StakingStatus int

const (
	// Unstaked represents unstaked status
	Unstaked StakingStatus = iota
	// Unstaking represents unstaking status
	Unstaking
	// Staked represents staked status
	Staked
)
