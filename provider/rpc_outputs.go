package provider

import (
	"fmt"
	"math/big"
	"time"
)

type queryBalanceOutput struct {
	Balance *big.Int `json:"balance"`
}

// GetAppOutput represents output for GetApp request
type GetAppOutput struct {
	Address       string    `json:"address"`
	PublicKey     string    `json:"public_key"`
	Jailed        bool      `json:"jailed"`
	Status        int       `json:"status"`
	Chains        []string  `json:"chains"`
	Tokens        string    `json:"tokens"`
	MaxRelays     string    `json:"max_relays"`
	UnstakingTime time.Time `json:"unstaking_time"`
}

// GetAppsOutput represents output for GetApps request
type GetAppsOutput struct {
	Result     []GetAppOutput `json:"result"`
	Page       int            `json:"page"`
	TotalPages int            `json:"total_pages"`
}

// GetNodeOutput represents output for GetNode request
type GetNodeOutput struct {
	Address       string    `json:"address"`
	Chains        []string  `json:"chains"`
	Jailed        bool      `json:"jailed"`
	PublicKey     string    `json:"public_key"`
	ServiceURL    string    `json:"service_url"`
	Status        int       `json:"status"`
	Tokens        string    `json:"tokens"`
	UnstakingTime time.Time `json:"unstaking_time"`
}

// GetNodesOutput represents output for GetNodes request
type GetNodesOutput struct {
	Result     []GetNodeOutput `json:"result"`
	Page       int             `json:"page"`
	TotalPages int             `json:"total_pages"`
}

// GetBlockOutput represents output for GetBlock request
type GetBlockOutput struct {
	Block struct {
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence []any `json:"evidence"`
		} `json:"evidence"`
		Header struct {
			AppHash       string `json:"app_hash"`
			ChainID       string `json:"chain_id"`
			ConsensusHash string `json:"consensus_hash"`
			DataHash      string `json:"data_hash"`
			EvidenceHash  string `json:"evidence_hash"`
			Height        string `json:"height"`
			LastBlockID   struct {
				Hash  string `json:"hash"`
				Parts struct {
					Hash  string `json:"hash"`
					Total string `json:"total"`
				} `json:"parts"`
			} `json:"last_block_id"`
			LastCommitHash     string    `json:"last_commit_hash"`
			LastResultsHash    string    `json:"last_results_hash"`
			NextValidatorsHash string    `json:"next_validators_hash"`
			NumTxs             string    `json:"num_txs"`
			ProposerAddress    string    `json:"proposer_address"`
			Time               time.Time `json:"time"`
			TotalTxs           string    `json:"total_txs"`
			ValidatorsHash     string    `json:"validators_hash"`
			Version            struct {
				App   string `json:"app"`
				Block string `json:"block"`
			} `json:"version"`
		} `json:"header"`
		LastCommit struct {
			BlockID struct {
				Hash  string `json:"hash"`
				Parts struct {
					Hash  string `json:"hash"`
					Total string `json:"total"`
				} `json:"parts"`
			} `json:"block_id"`
			Precommits []struct {
				BlockID struct {
					Hash  string `json:"hash"`
					Parts struct {
						Hash  string `json:"hash"`
						Total string `json:"total"`
					} `json:"parts"`
				} `json:"block_id"`
				Height           string    `json:"height"`
				Round            string    `json:"round"`
				Signature        string    `json:"signature"`
				Timestamp        time.Time `json:"timestamp"`
				Type             int       `json:"type"`
				ValidatorAddress string    `json:"validator_address"`
				ValidatorIndex   string    `json:"validator_index"`
			} `json:"precommits"`
		} `json:"last_commit"`
	} `json:"block"`
	BlockID struct {
		Hash  string `json:"hash"`
		Parts struct {
			Hash  string `json:"hash"`
			Total string `json:"total"`
		} `json:"parts"`
	} `json:"block_id"`
}

type queryHeightOutput struct {
	Height int `json:"height"`
}

// GetAccountOutput represents output for GetAccount request
type GetAccountOutput struct {
	Address string `json:"address"`
	Coins   []struct {
		Amount string `json:"amount"`
		Denom  string `json:"denom"`
	} `json:"coins"`
	PublicKey string `json:"public_key"`
}

// DispatchOutput represents output for Dispatch request
type DispatchOutput struct {
	BlockHeight int      `json:"block_height"`
	Session     *Session `json:"session"`
}

// Session represents session output from RPC request
type Session struct {
	Header *SessionHeader `json:"header"`
	Key    string         `json:"key"`
	Nodes  []*Node        `json:"nodes"`
}

// SessionHeader represents the headers of a session output
type SessionHeader struct {
	AppPublicKey  string `json:"app_public_key"`
	Chain         string `json:"chain"`
	SessionHeight int    `json:"session_height"`
}

// Node represents node output from RPC request
type Node struct {
	Address       string    `json:"address"`
	Chains        []string  `json:"chains"`
	Jailed        bool      `json:"jailed"`
	PublicKey     string    `json:"public_key"`
	ServiceURL    string    `json:"service_url"`
	Status        int       `json:"status"`
	Tokens        string    `json:"tokens"`
	UnstakingTime time.Time `json:"unstaking_time"`
}

// RPCError reprensents error output from RPC request
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error returns string representation of error
// needed to implement error interface
func (e *RPCError) Error() string {
	return fmt.Sprintf("Request failed with code: %v and message: %s", e.Code, e.Message)
}
