package provider

import (
	"github.com/pokt-foundation/pocket-go/pkg/models"
)

// GetNodeOptions represents optional arguments for GetNode request
type GetNodeOptions struct {
	Height int
}

// GetNodesOptions represents optional arguments for GetNodes request
type GetNodesOptions struct {
	StakingStatus models.StakingStatus
	Page          int
	PerPage       int
	Chain         string
	JailedStatus  models.JailedStatus
	Blockchain    string
}

// GetAppOptions represents optional arguments for GetApp request
type GetAppOptions struct {
	Height int
}

// GetAppsOptions represents optional arguments for GetApps request
type GetAppsOptions struct {
	StakingStatus models.StakingStatus
	Page          int
	PerPage       int
	Chain         string
	JailedStatus  models.JailedStatus
	Blockchain    string
}

// DispatchRequestOptions represents optional arguments for Dispatch request
type DispatchRequestOptions struct {
	RejectSelfSignedCertificates bool
}

// RelayRequestOptions represents optional arguments for Relay request
type RelayRequestOptions struct {
	RejectSelfSignedCertificates bool
}
