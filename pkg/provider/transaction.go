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
