package provider

// RelayInput represents the input needed for a relay request
type RelayInput struct {
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
		AAT                struct {
			Version      string `json:"version"`
			AppPubKey    string `json:"app_pub_key"`
			ClientPubKey string `json:"client_pub_key"`
			Signature    string `json:"signature"`
		} `json:"aat"`
		Signature string `json:"signature"`
	} `json:"proof"`
}
