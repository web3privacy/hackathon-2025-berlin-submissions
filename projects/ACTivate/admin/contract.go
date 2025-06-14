package admincontract

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethersphere/bee/v2/pkg/sctx"
	"github.com/ethersphere/bee/v2/pkg/transaction"
)

var (
	ErrBagoy         = errors.New("")
	bagoydescription = "bagoy"
)

type Interface interface {
	SendDataToTarget(ctx context.Context, target common.Address, owner, actRef []byte, topic string) (receipt *types.Receipt, err error)
}

type admincontract struct {
	owner                common.Address
	adminContractAddress common.Address // 0x8946CCb5176E614F342157139c7546DA81085E6f
	adminContractABI     abi.ABI
	transactionService   transaction.Service
	gasLimit             uint64
	dataSentToTarget     common.Hash
}

func New(
	owner common.Address,
	adminContractAddress common.Address,
	adminContractABI abi.ABI,
	transactionService transaction.Service,
	setGasLimit bool,
) Interface {

	var gasLimit uint64
	if setGasLimit {
		gasLimit = transaction.DefaultGasLimit
	}

	return &admincontract{
		owner:                owner,
		adminContractAddress: adminContractAddress,
		adminContractABI:     adminContractABI,
		transactionService:   transactionService,
		gasLimit:             gasLimit,
		dataSentToTarget:     adminContractABI.Events["DataSentToTarget"].ID,
	}
}

func (c *admincontract) SendDataToTarget(ctx context.Context, target common.Address, owner, actRef []byte, topic string) (receipt *types.Receipt, err error) {

	callData, err := c.adminContractABI.Pack("sendDataToTarget", target, owner, actRef, topic)
	if err != nil {
		return nil, err
	}

	receipt, err = c.sendTransaction(ctx, callData, "sendDataToTarget")
	if err != nil {
		return nil, fmt.Errorf("send data to target: %w", err)
	}

	return receipt, nil
}

func (c *admincontract) sendTransaction(ctx context.Context, callData []byte, desc string) (receipt *types.Receipt, err error) {
	request := &transaction.TxRequest{
		To:          &c.adminContractAddress,
		Data:        callData,
		GasPrice:    sctx.GetGasPrice(ctx),
		GasLimit:    max(sctx.GetGasLimit(ctx), c.gasLimit),
		Value:       big.NewInt(0),
		Description: desc,
	}

	defer func() {
		err = c.transactionService.UnwrapABIError(
			ctx,
			request,
			err,
			c.adminContractABI.Errors,
		)
	}()

	txHash, err := c.transactionService.Send(ctx, request, transaction.DefaultTipBoostPercent)
	if err != nil {
		return nil, err
	}

	receipt, err = c.transactionService.WaitForReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}

	if receipt.Status == 0 {
		return nil, transaction.ErrTransactionReverted
	}

	return receipt, nil
}
