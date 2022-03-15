package provider

import (
	"math/big"
	"time"
)

type queryBalanceResponse struct {
	Balance *big.Int `json:"balance"`
}

type queryAccountsTXsResponse struct {
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
			Fee     struct {
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
	TotalCount *string `json:"total_count"`
}

// GetAppResponse represents response for GetApp request
type GetAppResponse struct {
	Address       string    `json:"address"`
	PublicKey     string    `json:"public_key"`
	Jailed        bool      `json:"jailed"`
	Status        int       `json:"status"`
	Chains        []string  `json:"chains"`
	Tokens        string    `json:"tokens"`
	MaxRelays     *int      `json:"max_relays"`
	UnstakingTime time.Time `json:"unstaking_time"`
}

// GetAppsResponse represents response for GetApps request
type GetAppsResponse struct {
	Result     []GetAppResponse `json:"result"`
	Page       int              `json:"page"`
	TotalPages int              `json:"total_pages"`
}

// GetNodeResponse represents response for GetNode request
type GetNodeResponse struct {
	Address       string    `json:"address"`
	Chains        *[]string `json:"chains"`
	Jailed        bool      `json:"jailed"`
	PublicKey     string    `json:"public_key"`
	ServiceURL    *string   `json:"service_url"`
	Status        int       `json:"status"`
	Tokens        string    `json:"tokens"`
	UnstakingTime time.Time `json:"unstaking_time"`
}

// GetNodesResponse represents response for GetNodes request
type GetNodesResponse struct {
	Result     []GetNodeResponse `json:"result"`
	Page       int               `json:"page"`
	TotalPages int               `json:"total_pages"`
}

// SendTransactionResponse represents response for SendTransaction request
type SendTransactionResponse struct {
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

// GetBlockResponse represents response for GetBlock request
type GetBlockResponse struct {
	Block *struct {
		Data struct {
			Txs string `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence string `json:"evidence"`
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
			Precommits interface{} `json:"precommits"`
		} `json:"last_commit"`
	} `json:"block"`
	BlockMeta struct {
		BlockID struct {
			Hash  string `json:"hash"`
			Parts struct {
				Hash  string `json:"hash"`
				Total string `json:"total"`
			} `json:"parts"`
		} `json:"block_id"`
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
	} `json:"block_meta"`
}

// GetTransactionResponse represents response for GetTransaction request
type GetTransactionResponse struct {
	Transaction struct {
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
			Fee     struct {
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
	} `json:"transaction"`
}

type queryHeightResponse struct {
	Height *int `json:"height"`
}

// GetAccountResponse represents response for GetAccount request
type GetAccountResponse struct {
	Address string `json:"address"`
	Coins   []struct {
		Amount string `json:"amount"`
		Denom  string `json:"denom"`
	} `json:"coins"`
	PublicKey string `json:"public_key"`
}

// GetAccountWithTransactionsResponse represents response for GetAccountWithTransactions request
type GetAccountWithTransactionsResponse struct {
	Account      *GetAccountResponse
	Transactions *queryAccountsTXsResponse
}

// DispatchResponse represents response for Dispatch request
type DispatchResponse struct {
	BlockHeight int `json:"block_height"`
	Session     *struct {
		Header struct {
			AppPublicKey  string `json:"app_public_key"`
			Chain         string `json:"chain"`
			SessionHeight int    `json:"session_height"`
		} `json:"header"`
		Key   string `json:"key"`
		Nodes []struct {
			Address       string    `json:"address"`
			Chains        []string  `json:"chains"`
			Jailed        bool      `json:"jailed"`
			PublicKey     string    `json:"public_key"`
			ServiceURL    string    `json:"service_url"`
			Status        int       `json:"status"`
			Tokens        string    `json:"tokens"`
			UnstakingTime time.Time `json:"unstaking_time"`
		} `json:"nodes"`
	} `json:"session"`
}

// RelayResponse represents response for Relay request
type RelayResponse struct {
	Payload struct {
		Data   string `json:"data"`
		Method string `json:"method"`
		Path   string `json:"path"`
	} `json:"payload"`
	Meta struct {
		BlockHeight int `json:"block_height"`
	} `json:"meta"`
	Proof struct {
		RequestHash        string `json:"request_hash"`
		Entropy            int64  `json:"entropy"`
		SessionBlockHeight int    `json:"session_block_height"`
		ServicerPubKey     string `json:"servicer_pub_key"`
		Blockchain         string `json:"blockchain"`
		Aat                struct {
			Version      string `json:"version"`
			AppPubKey    string `json:"app_pub_key"`
			ClientPubKey string `json:"client_pub_key"`
			Signature    string `json:"signature"`
		} `json:"aat"`
		Signature string `json:"signature"`
	} `json:"proof"`
}
