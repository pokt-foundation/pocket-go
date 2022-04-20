package transactionbuilder

import (
	"encoding/hex"

	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	appsType "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

// NewMsgSend returns message for send transaction
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

// NewStakeApp returns message for Stake App transaction
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

// NewUnstakeApp returns message for Unstake App transaction
func NewUnstakeApp(address string) (TxMsg, error) {
	decodedAddress, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return &appsType.MsgBeginUnstake{
		Address: decodedAddress,
	}, nil
}

// NewUnjailApp returns message for Unjail App transaction
func NewUnjailApp(address string) (TxMsg, error) {
	decodedAddress, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return &appsType.MsgUnjail{
		AppAddr: decodedAddress,
	}, nil
}

// NewStakeNode returns message for Stake Node transaction
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

// NewUnstakeNode returns message for Unstake Node transaction
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

// NewUnjailNode returns message for Unjail Node transaction
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
