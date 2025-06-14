package screens

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum" // Required for FilterQuery and Subscription
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient" // Required for ethclient.Client
	"github.com/ethersphere/bee/v2/pkg/sctx"
	"github.com/ethersphere/bee/v2/pkg/transaction"
)

var (
	ErrBagoy         = errors.New("")
	bagoydescription = "bagoy"
)

type DataContractInterface interface {
	SendDataToTarget(ctx context.Context, target common.Address, owner, actRef []byte, topic string) (receipt *types.Receipt, err error)
	SubscribeDataSentToTarget(ctx context.Context, client *ethclient.Client, sink chan<- types.Log) (ethereum.Subscription, error)
}

type datacontract struct {
	owner               common.Address
	dataContractAddress common.Address
	dataContractABI     abi.ABI
	transactionService  transaction.Service
	// ethClient field is not strictly needed here if passed directly to SubscribeDataSentToTarget
	gasLimit         uint64
	dataSentToTarget common.Hash // This is the event signature HASH for DataSentToTarget
}

func NewDataContract(
	owner common.Address,
	dataContractAddress common.Address,
	dataContractABI abi.ABI,
	transactionService transaction.Service,
	// ethClient *ethclient.Client, // Not adding to struct for now, passed directly to subscribe method
	setGasLimit bool,
) DataContractInterface {

	var gasLimit uint64
	if setGasLimit {
		gasLimit = transaction.DefaultGasLimit
	}

	return &datacontract{
		owner:               owner,
		dataContractAddress: dataContractAddress,
		dataContractABI:     dataContractABI,
		transactionService:  transactionService,
		gasLimit:            gasLimit,
		dataSentToTarget:    dataContractABI.Events["DataSentToTarget"].ID,
	}
}

func (c *datacontract) SendDataToTarget(ctx context.Context, target common.Address, owner, actRef []byte, topic string) (receipt *types.Receipt, err error) {

	callData, err := c.dataContractABI.Pack("sendDataToTarget", target, owner, actRef, topic)
	if err != nil {
		return nil, err
	}

	receipt, err = c.sendTransaction(ctx, callData, "sendDataToTarget")
	if err != nil {
		return nil, fmt.Errorf("send data to target: %w", err)
	}

	return receipt, nil
}

// rpcClient, err := rpc.DialContext(ctx, endpoint)
// ethclient.NewClient(rpcClient)
func (c *datacontract) SubscribeDataSentToTarget(ctx context.Context, client *ethclient.Client, sink chan<- types.Log) (ethereum.Subscription, error) {
	if client == nil {
		return nil, errors.New("ethclient.Client is nil")
	}

	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("Error getting current block number for admin subscription:", err)
	} else {
		currentBlock = 40581246
	}
	log.Printf("Obtained currentBlock: %d for admin subscription", currentBlock)

	fromBlockBigInt := new(big.Int).SetUint64(currentBlock)

	log.Printf("Subscribing to DataContract DataSentToTarget events from block %s", fromBlockBigInt.String())

	const blockPageSize = 500
	query := ethereum.FilterQuery{
		Addresses: []common.Address{c.dataContractAddress},
		Topics:    [][]common.Hash{{c.dataSentToTarget}},
		FromBlock: fromBlockBigInt,
		ToBlock:   big.NewInt(int64(currentBlock + blockPageSize - 1)),
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to DataSentToTarget events: %w", err)
	}

	// Start a goroutine to forward logs from the subscription channel to the sink channel
	// This goroutine also handles unsubscription and error propagation
	go func() {
		defer close(sink) // Close the sink channel when the goroutine exits
		for {
			select {
			case <-ctx.Done(): // If the context is cancelled, stop the goroutine
				fmt.Println("Subscription context cancelled, unsubscribing.")
				sub.Unsubscribe()
				return
			case err := <-sub.Err(): // If the subscription encounters an error
				// Log the error. Depending on the error, you might want to attempt to resubscribe
				// or signal the calling code that the subscription has failed.
				fmt.Printf("Event subscription error for DataSentToTarget: %v\\n", err)
				// Forwarding the error or handling reconnection is an advanced topic.
				// For now, we just stop this goroutine.
				return
			case vLog := <-logs: // Received a new log
				// Send the raw log to the sink channel.
				// The consumer of the sink channel will be responsible for unpacking the log
				// using c.dataContractABI.UnpackIntoInterface or similar.
				select {
				case sink <- vLog:
				case <-ctx.Done(): // Check context again before sending to avoid blocking
					fmt.Println("Subscription context cancelled while trying to send log, unsubscribing.")
					sub.Unsubscribe()
					return
				}
			}
		}
	}()

	return sub, nil
}

func (c *datacontract) sendTransaction(ctx context.Context, callData []byte, desc string) (receipt *types.Receipt, err error) {
	request := &transaction.TxRequest{
		To:          &c.dataContractAddress,
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
			c.dataContractABI.Errors,
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
