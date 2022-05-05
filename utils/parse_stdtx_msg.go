package utils

import (
	transactionbuilder "github.com/pokt-foundation/pocket-go/transaction-builder"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

type StdTxMsg struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func ParseStdTxMsg(msg StdTxMsg) (transactionbuilder.TransactionMessage, error) {
	switch msg.Type {
	case "pos/Send":
		parsedValue := msg.Value.(*nodesTypes.MsgSend)
		return transactionbuilder.NewSend(parsedValue.FromAddress.String(), parsedValue.ToAddress.String(), parsedValue.Amount.Int64())
	default:
		return nil, nil
	}
}
