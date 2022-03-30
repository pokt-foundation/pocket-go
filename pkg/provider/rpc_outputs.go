package provider

import (
	"fmt"
	"math/big"
	"time"
)

type queryBalanceOutput struct {
	Balance *big.Int `json:"balance"`
}

// GetAccountTransactionsOutput represents output for GetAccountTransactions request
type GetAccountTransactionsOutput struct {
	Txs []struct {
		Hash     string `json:"hash"`
		Height   int    `json:"height"`
		Index    int    `json:"index"`
		TxResult struct {
			Code        int      `json:"code"`
			Data        string   `json:"data"`
			Log         string   `json:"log"`
			Info        string   `json:"info"`
			Events      []string `json:"events"`
			Codespace   string   `json:"codespace"`
			Signer      string   `json:"signer"`
			Recipient   string   `json:"recipient"`
			MessageType string   `json:"message_type"`
		} `json:"tx_result"`
		Tx    string `json:"tx"`
		Proof struct {
			RootHash string `json:"root_hash"`
			Data     string `json:"data"`
			Proof    struct {
				Total    int      `json:"total"`
				Index    int      `json:"index"`
				LeafHash string   `json:"leaf_hash"`
				Aunts    []string `json:"aunts"`
			} `json:"proof"`
		} `json:"proof"`
		StdTx struct {
			Entropy int `json:"entropy"`
			Fee     []struct {
				Amount string `json:"amount"`
				Denom  string `json:"denom"`
			} `json:"fee"`
			Memo string `json:"memo"`
			Msg  struct {
			} `json:"msg"`
			Signature struct {
				PubKey    string `json:"pub_key"`
				Signature string `json:"signature"`
			} `json:"signature"`
		} `json:"stdTx"`
	} `json:"txs"`
	TotalCount int `json:"total_count"`
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

// SendTransactionOutput represents output for SendTransaction request
type SendTransactionOutput struct {
	Height string `json:"height"`
	Txhash string `json:"txhash"`
	RawLog string `json:"raw_log"`
	Logs   []struct {
		MsgIndex int    `json:"msg_index"`
		Success  bool   `json:"success"`
		Log      string `json:"log"`
		Events   []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"logs"`
}

// GetBlockOutput represents output for GetBlock request
type GetBlockOutput struct {
	Block struct {
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence interface{} `json:"evidence"`
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
			Precommits []interface{} `json:"precommits"`
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

// GetTransactionOutput represents output for GetTransaction request
type GetTransactionOutput struct {
	Hash     string `json:"hash"`
	Height   int    `json:"height"`
	Index    int    `json:"index"`
	TxResult struct {
		Code        int      `json:"code"`
		Data        string   `json:"data"`
		Log         string   `json:"log"`
		Info        string   `json:"info"`
		Events      []string `json:"events"`
		Codespace   string   `json:"codespace"`
		Signer      string   `json:"signer"`
		Recipient   string   `json:"recipient"`
		MessageType string   `json:"message_type"`
	} `json:"tx_result"`
	Tx    string `json:"tx"`
	Proof struct {
		RootHash string `json:"root_hash"`
		Data     string `json:"data"`
		Proof    struct {
			Total    int      `json:"total"`
			Index    int      `json:"index"`
			LeafHash string   `json:"leaf_hash"`
			Aunts    []string `json:"aunts"`
		} `json:"proof"`
	} `json:"proof"`
	StdTx struct {
		Entropy int `json:"entropy"`
		Fee     []struct {
			Amount string `json:"amount"`
			Denom  string `json:"denom"`
		} `json:"fee"`
		Memo string `json:"memo"`
		Msg  struct {
		} `json:"msg"`
		Signature struct {
			PubKey    string `json:"pub_key"`
			Signature string `json:"signature"`
		} `json:"signature"`
	} `json:"stdTx"`
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
