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
	TotalCount string `json:"total_count"`
}

type queryAppResponse struct {
	Address       string    `json:"address"`
	PublicKey     string    `json:"public_key"`
	Jailed        bool      `json:"jailed"`
	Status        int       `json:"status"`
	Chains        []string  `json:"chains"`
	Tokens        string    `json:"tokens"`
	MaxRelays     *int      `json:"max_relays"`
	UnstakingTime time.Time `json:"unstaking_time"`
}

type queryNodeResponse struct {
	Address       string    `json:"address"`
	Chains        []string  `json:"chains"`
	Jailed        bool      `json:"jailed"`
	PublicKey     string    `json:"public_key"`
	ServiceURL    *string   `json:"service_url"`
	Status        int       `json:"status"`
	Tokens        string    `json:"tokens"`
	UnstakingTime time.Time `json:"unstaking_time"`
}

type TransactionReponse struct {
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
