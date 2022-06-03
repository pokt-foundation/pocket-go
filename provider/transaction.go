package provider

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

// SendTransactionInput represents input needed for SendTransaction request
type SendTransactionInput struct {
	Address     string `json:"address"`
	RawHexBytes string `json:"raw_hex_bytes"`
}

// GetTransactionOutput represents output for GetTransaction request
type GetTransactionOutput struct {
	*Transaction
}

// GetAccountTransactionsOutput represents output for GetAccountTransactions request
type GetAccountTransactionsOutput struct {
	PageCount int            `json:"page_count"`
	TotalTxs  int            `json:"total_txs"`
	Txs       []*Transaction `json:"txs"`
}

// GetBlockTransactionsOutput represents output for GetBlockTransactions request
type GetBlockTransactionsOutput struct {
	PageCount int            `json:"page_count"`
	TotalTxs  int            `json:"total_txs"`
	Txs       []*Transaction `json:"txs"`
}

// TransactionProof represents the proof of transaction
type TransactionProof struct {
	Data  string `json:"data"`
	Proof struct {
		Aunts    []string `json:"aunts"`
		Index    int      `json:"index"`
		LeafHash string   `json:"leaf_hash"`
		Total    int      `json:"total"`
	} `json:"proof"`
	RootHash string `json:"root_hash"`
}

// StdTx represents standard transaction fields
type StdTx struct {
	Entropy   int64        `json:"entropy"`
	Fee       []*Fee       `json:"fee"`
	Memo      string       `json:"memo"`
	Msg       *TxMsg       `json:"msg"`
	Signature *TxSignature `json:"signature"`
}

// Fee represents fee values
type Fee struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

// TxMsg represents message for transactions
type TxMsg struct {
	Type  string         `json:"type"`
	Value map[string]any `json:"value"`
}

// TxSignature represents values of transaction signature
type TxSignature struct {
	PubKey    string `json:"pub_key"`
	Signature string `json:"signature"`
}

// TxResult represents transaction result
type TxResult struct {
	Code        int    `json:"code"`
	Codespace   string `json:"codespace"`
	Data        string `json:"data"`
	Events      string `json:"events"`
	Info        string `json:"info"`
	Log         string `json:"log"`
	MessageType string `json:"message_type"`
	Recipient   string `json:"recipient"`
	Signer      string `json:"signer"`
}

// Transaction represents a transaction in Pocket
type Transaction struct {
	Hash     string            `json:"hash"`
	Height   int               `json:"height"`
	Index    int               `json:"index"`
	Proof    *TransactionProof `json:"proof"`
	StdTx    *StdTx            `json:"stdTx"`
	Tx       string            `json:"tx"`
	TxResult *TxResult         `json:"tx_result"`
}
