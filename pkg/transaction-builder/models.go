package transactionbuilder

import (
	"encoding/hex"

	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	appsType "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

func NewMsgSend(fromAddress, toAddress string, amount int64) (TxMsg, error) {
	decodedFromAddress, err := hex.DecodeString(fromAddress)
	if err != nil {
		return nil, err
	}

	decodedToAddress, err := hex.DecodeString(toAddress)
	if err != nil {
		return nil, err
	}

	return &nodeTypes.MsgSend{
		FromAddress: decodedFromAddress,
		ToAddress:   decodedToAddress,
		Amount:      types.NewInt(amount),
	}, nil
}

func NewStakeApp(publicKey string, chains []string, amount int64) (TxMsg, error) {
	cryptoPublicKey, err := crypto.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &appsType.MsgStake{
		PubKey: cryptoPublicKey,
		Chains: chains,
		Value:  types.NewInt(amount),
	}, nil
}

func NewUnstakeApp(address string) (TxMsg, error) {
	decodedAddress, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return &appsType.MsgBeginUnstake{
		Address: decodedAddress,
	}, nil
}

func NewUnjailApp(address string) (TxMsg, error) {
	decodedAddress, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return &appsType.MsgUnjail{
		AppAddr: decodedAddress,
	}, nil
}

func NewStakeNode(publicKey, serviceUrl, outputAddress string, chains []string, amount int64) (TxMsg, error) {
	cryptoPublicKey, err := crypto.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	decodedAddress, err := hex.DecodeString(outputAddress)
	if err != nil {
		return nil, err
	}

	return &nodeTypes.MsgStake{
		PublicKey:  cryptoPublicKey,
		Chains:     chains,
		Value:      types.NewInt(amount),
		ServiceUrl: serviceUrl,
		Output:     decodedAddress,
	}, nil
}

func NewUnstakeNode(fromAddress, operatorAddress string) (TxMsg, error) {
	decodedFromAddress, err := hex.DecodeString(fromAddress)
	if err != nil {
		return nil, err
	}

	decodedOperatorAddress, err := hex.DecodeString(operatorAddress)
	if err != nil {
		return nil, err
	}

	return &nodeTypes.MsgBeginUnstake{
		Address: decodedOperatorAddress,
		Signer:  decodedFromAddress,
	}, nil
}

func NewUnjailNode(fromAddress, operatorAddress string) (TxMsg, error) {
	decodedFromAddress, err := hex.DecodeString(fromAddress)
	if err != nil {
		return nil, err
	}

	decodedOperatorAddress, err := hex.DecodeString(operatorAddress)
	if err != nil {
		return nil, err
	}

	return &nodeTypes.MsgUnjail{
		ValidatorAddr: decodedOperatorAddress,
		Signer:        decodedFromAddress,
	}, nil
}
