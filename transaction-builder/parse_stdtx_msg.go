package transactionbuilder

import (
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

// StdTxMsg represents 'StdTx' field from a Transaction
type StdTxMsg struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// ParseStdTxMsg parses any StdTx.Msg.Value and returns a typed struct
func ParseStdTxMsg(msg StdTxMsg) (TransactionMessage, error) {
	switch msg.Type {
	case "pos/Send":
		parsedValue := msg.Value.(*nodesTypes.MsgSend)
		return NewSend(parsedValue.FromAddress.String(), parsedValue.ToAddress.String(), parsedValue.Amount.Int64())
	default:
		return nil, nil
	}
}
