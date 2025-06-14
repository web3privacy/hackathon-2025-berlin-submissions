package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// AdminContract ABI - Generated from compiled contract
const AdminContractABI = `[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "from",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "to",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "bytes32",
				"name": "owner",
				"type": "bytes32"
			},
			{
				"indexed": false,
				"internalType": "bytes32",
				"name": "actref",
				"type": "bytes32"
			},
			{
				"indexed": false,
				"internalType": "string",
				"name": "topic",
				"type": "string"
			}
		],
		"name": "DataSentToTarget",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "getAdmin",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "target",
				"type": "address"
			},
			{
				"internalType": "bytes32",
				"name": "ownerParam",
				"type": "bytes32"
			},
			{
				"internalType": "bytes32",
				"name": "actref",
				"type": "bytes32"
			},
			{
				"internalType": "string",
				"name": "topic",
				"type": "string"
			}
		],
		"name": "sendDataToTarget",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

// DataSentToTargetEvent represents the event emitted by the contract
type DataSentToTargetEvent struct {
	From   common.Address
	To     common.Address
	Owner  [32]byte
	Actref [32]byte
	Topic  string
}

func main() {
	fmt.Println("=== AdminContract Go Interaction Example ===\n")

	// Get environment variables
	rpcURL := getEnv("RPC_URL", "http://localhost:8545")
	contractAddressHex := getEnv("CONTRACT_ADDRESS", "")
	privateKeyHex := getEnv("PRIVATE_KEY", "")

	if contractAddressHex == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable not set")
	}
	if privateKeyHex == "" {
		log.Fatal("PRIVATE_KEY environment variable not set")
	}

	// Connect to Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	fmt.Printf("Connected to: %s\n", rpcURL)

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Using account: %s\n", fromAddress.Hex())

	// Parse contract address
	contractAddress := common.HexToAddress(contractAddressHex)
	fmt.Printf("Contract address: %s\n", contractAddress.Hex())

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(AdminContractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Create contract instance
	contract := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)

	// Get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(300000) // Set gas limit

	fmt.Println("\n=== Contract Information ===")

	// Get admin address
	var adminResult []interface{}
	err = contract.Call(&bind.CallOpts{}, &adminResult, "getAdmin")
	if err != nil {
		log.Fatalf("Failed to get admin: %v", err)
	}
	adminAddress := adminResult[0].(common.Address)
	fmt.Printf("Admin address: %s\n", adminAddress.Hex())
	fmt.Printf("Is caller admin? %t\n", adminAddress == fromAddress)

	if adminAddress != fromAddress {
		fmt.Println("WARNING: Connected account is not the admin. Transaction will fail.")
	}

	fmt.Println("\n=== Sending Data to Target ===")

	// Example target address (could be any address)
	targetAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b8d6968e2a4aF9c11B")

	// Convert strings to bytes32
	ownerParam := stringToBytes32("OWNER_GO_001")
	actrefParam := stringToBytes32("ACTREF_GO_123")
	topic := "Go Example - Blockchain Integration"

	fmt.Printf("Target address: %s\n", targetAddress.Hex())
	fmt.Printf("Owner param: %s\n", fmt.Sprintf("0x%x", ownerParam))
	fmt.Printf("Action ref: %s\n", fmt.Sprintf("0x%x", actrefParam))
	fmt.Printf("Topic: %s\n", topic)

	// Send transaction
	fmt.Println("\nSending transaction...")
	tx, err := contract.Transact(auth, "sendDataToTarget", targetAddress, ownerParam, actrefParam, topic)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
	fmt.Println("Waiting for confirmation...")

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction: %v", err)
	}

	fmt.Printf("Transaction confirmed in block: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas used: %d\n", receipt.GasUsed)

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Println("✓ Transaction successful!")
	} else {
		fmt.Println("✗ Transaction failed!")
	}

	// Parse events from receipt
	fmt.Println("\n=== Event Details ===")
	for _, vLog := range receipt.Logs {
		if vLog.Address == contractAddress && len(vLog.Topics) > 0 {
			// Check if this is our DataSentToTarget event
			eventSignature := vLog.Topics[0]
			expectedSignature := crypto.Keccak256Hash([]byte("DataSentToTarget(address,address,bytes32,bytes32,string)"))

			if eventSignature == expectedSignature {
				fmt.Println("Found DataSentToTarget event:")

				// Parse indexed parameters from topics
				if len(vLog.Topics) >= 3 {
					from := common.HexToAddress(vLog.Topics[1].Hex())
					to := common.HexToAddress(vLog.Topics[2].Hex())
					fmt.Printf("  From: %s\n", from.Hex())
					fmt.Printf("  To: %s\n", to.Hex())
				}

				// Parse non-indexed parameters from data
				if len(vLog.Data) > 0 {
					var event DataSentToTargetEvent
					err := parsedABI.UnpackIntoInterface(&event, "DataSentToTarget", vLog.Data)
					if err != nil {
						fmt.Printf("  Error parsing event data: %v\n", err)
					} else {
						fmt.Printf("  Owner: 0x%x\n", event.Owner)
						fmt.Printf("  Action Ref: 0x%x\n", event.Actref)
						fmt.Printf("  Topic: %s\n", event.Topic)
					}
				}
			}
		}
	}

	fmt.Println("\n=== Reading Contract State ===")

	// Query contract state
	var ownerResult []interface{}
	err = contract.Call(&bind.CallOpts{}, &ownerResult, "owner")
	if err != nil {
		log.Printf("Failed to call owner(): %v", err)
	} else {
		ownerAddr := ownerResult[0].(common.Address)
		fmt.Printf("Contract owner: %s\n", ownerAddr.Hex())
	}

	fmt.Println("\n=== Event Filtering Example ===")

	// Filter events from the last 10 blocks
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Printf("Failed to get current block: %v", err)
	} else {
		fromBlock := currentBlock - 10
		if fromBlock > currentBlock { // Handle underflow
			fromBlock = 0
		}

		fmt.Printf("Filtering events from block %d to %d...\n", fromBlock, currentBlock)

		// Create filter query
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(currentBlock)),
			Addresses: []common.Address{contractAddress},
			Topics: [][]common.Hash{
				{crypto.Keccak256Hash([]byte("DataSentToTarget(address,address,bytes32,bytes32,string)"))},
			},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("Failed to filter logs: %v", err)
		} else {
			fmt.Printf("Found %d events:\n", len(logs))
			for i, vLog := range logs {
				fmt.Printf("  Event %d - Block: %d, TxHash: %s\n",
					i+1, vLog.BlockNumber, vLog.TxHash.Hex())
			}
		}
	}

	fmt.Println("\n=== Go Integration Complete ===")
	fmt.Println("Successfully demonstrated AdminContract interaction from Go!")
	fmt.Println("\nTo run this example:")
	fmt.Println("1. Set environment variables:")
	fmt.Println("   export RPC_URL=http://localhost:8545")
	fmt.Println("   export CONTRACT_ADDRESS=0x...")
	fmt.Println("   export PRIVATE_KEY=0x...")
	fmt.Println("2. Install dependencies: go mod init && go mod tidy")
	fmt.Println("3. Run: go run main.go")
}

// Helper function to convert string to bytes32
func stringToBytes32(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
