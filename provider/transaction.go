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

type Transaction struct {
	Hash   string `json:"hash"`
	Height int    `json:"height"`
	Index  int    `json:"index"`
	Proof  struct {
		Data  interface{} `json:"data"`
		Proof struct {
			Aunts    interface{} `json:"aunts"`
			Index    int         `json:"index"`
			LeafHash interface{} `json:"leaf_hash"`
			Total    int         `json:"total"`
		} `json:"proof"`
		RootHash string `json:"root_hash"`
	} `json:"proof"`
	StdTx struct {
		Entropy int64 `json:"entropy"`
		Fee     []struct {
			Amount string `json:"amount"`
			Denom  string `json:"denom"`
		} `json:"fee"`
		Memo string `json:"memo"`
		Msg  struct {
			Type  string `json:"type"`
			Value struct {
				Amount      string `json:"amount"`
				FromAddress string `json:"from_address"`
				ToAddress   string `json:"to_address"`
			} `json:"value"`
		} `json:"msg"`
		Signature struct {
			PubKey    string `json:"pub_key"`
			Signature string `json:"signature"`
		} `json:"signature"`
	} `json:"stdTx"`
	Tx       string `json:"tx"`
	TxResult struct {
		Code        int         `json:"code"`
		Codespace   string      `json:"codespace"`
		Data        interface{} `json:"data"`
		Events      interface{} `json:"events"`
		Info        string      `json:"info"`
		Log         string      `json:"log"`
		MessageType string      `json:"message_type"`
		Recipient   string      `json:"recipient"`
		Signer      string      `json:"signer"`
	} `json:"tx_result"`
}
