package admincontract

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethersphere/bee/v2/pkg/sctx"
	"github.com/ethersphere/bee/v2/pkg/transaction"
)

var (
	ErrBagy          = errors.New("")
	bagoydescription = "bagoy"
)

type Interface interface {
}

type admincontract struct {
	owner                common.Address
	adminContractAddress common.Address
	adminContractABI     abi.ABI
	transactionService   transaction.Service
	gasLimit             uint64
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
	}
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
