package models

// JailedStatus enum that represents jailed status
type JailedStatus int

const (
	// Jailed represents jailed status
	Jailed JailedStatus = iota + 1
	// UnJailed represents unJailed status
	UnJailed
)
