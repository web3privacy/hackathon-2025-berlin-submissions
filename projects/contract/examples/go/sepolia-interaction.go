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

// LoadDeploymentConfig loads the deployment configuration from environment or constants
func LoadDeploymentConfig() (string, string, string) {
	// Try to load from environment first
	rpcURL := os.Getenv("RPC_URL")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")

	// Fallback to constants (replace these with your actual deployment values)
	if rpcURL == "" {
		rpcURL = "https://sepolia.infura.io/v3/YOUR_PROJECT_ID" // Replace with actual RPC URL
	}
	if contractAddress == "" {
		contractAddress = "0x0000000000000000000000000000000000000000" // Replace with actual contract address
	}
	if privateKey == "" {
		privateKey = "0x0000000000000000000000000000000000000000000000000000000000000000" // Replace with actual private key
	}

	return rpcURL, contractAddress, privateKey
}

func main() {
	fmt.Println("=== AdminContract Sepolia Interaction ===\n")

	// Load deployment configuration
	rpcURL, contractAddressHex, privateKeyHex := LoadDeploymentConfig()

	if contractAddressHex == "0x0000000000000000000000000000000000000000" {
		log.Fatal(`❌ Contract address not set!
		
Please either:
1. Set environment variables:
   export RPC_URL="https://sepolia.infura.io/v3/YOUR_PROJECT_ID"
   export CONTRACT_ADDRESS="0x..."
   export PRIVATE_KEY="0x..."

2. Or update the constants in LoadDeploymentConfig() function with your deployment values.
   You can find these values in the deployments/constants.go file after deployment.`)
	}

	fmt.Printf("🔗 RPC URL: %s\n", rpcURL)
	fmt.Printf("📍 Contract Address: %s\n", contractAddressHex)

	// Connect to Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	// Verify we're connected to Sepolia
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	if chainID.Int64() != 11155111 {
		log.Printf("⚠️  Warning: Connected to chain ID %d, expected Sepolia (11155111)", chainID.Int64())
	} else {
		fmt.Println("✅ Connected to Sepolia testnet")
	}

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
	fmt.Printf("👤 Using account: %s\n", fromAddress.Hex())

	// Check account balance
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	balanceEth := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(big.NewInt(1e18)))
	fmt.Printf("💰 Account balance: %s ETH\n", balanceEth.String())

	// Parse contract address
	contractAddress := common.HexToAddress(contractAddressHex)

	// AdminContract ABI - This should match your deployed contract
	const AdminContractABI = `[
		{
			"inputs": [],
			"stateMutability": "nonpayable",
			"type": "constructor"
		},
		{
			"anonymous": false,
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

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(AdminContractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Create contract instance
	contract := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)

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
	auth.GasLimit = uint64(300000)

	fmt.Println("\n=== Contract Information ===")

	// Check if account is admin
	var adminResult []interface{}
	err = contract.Call(&bind.CallOpts{}, &adminResult, "owner")
	if err != nil {
		log.Printf("Failed to get contract admin: %v", err)
	} else {
		adminAddress := adminResult[0].(common.Address)
		fmt.Printf("📋 Contract admin: %s\n", adminAddress.Hex())
		fmt.Printf("🔑 Is caller admin? %t\n", adminAddress == fromAddress)

		if adminAddress != fromAddress {
			fmt.Println("⚠️  Warning: Connected account is not the admin. Transaction will fail.")
		}
	}

	fmt.Println("\n=== Sending Data to Target ===")

	// Example target address
	targetAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b8d6968e2a4aF9c11B")

	// Convert strings to bytes32
	ownerParam := stringToBytes32("SEPOLIA_OWNER_001")
	actrefParam := stringToBytes32("SEPOLIA_ACTION_123")
	topic := "Sepolia Testnet - Go Integration Test"

	fmt.Printf("🎯 Target address: %s\n", targetAddress.Hex())
	fmt.Printf("👤 Owner param: %s\n", fmt.Sprintf("0x%x", ownerParam))
	fmt.Printf("📝 Action ref: %s\n", fmt.Sprintf("0x%x", actrefParam))
	fmt.Printf("🏷️  Topic: %s\n", topic)

	// Send transaction
	fmt.Println("\n🚀 Sending transaction...")
	tx, err := contract.Transact(auth, "sendDataToTarget", targetAddress, ownerParam, actrefParam, topic)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("📋 Transaction hash: %s\n", tx.Hash().Hex())
	fmt.Printf("🔗 View on Etherscan: https://sepolia.etherscan.io/tx/%s\n", tx.Hash().Hex())
	fmt.Println("⏳ Waiting for confirmation...")

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction: %v", err)
	}

	fmt.Printf("✅ Transaction confirmed in block: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("⛽ Gas used: %d\n", receipt.GasUsed)

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Println("🎉 Transaction successful!")
	} else {
		fmt.Println("❌ Transaction failed!")
	}

	// Parse events from receipt
	fmt.Println("\n=== Event Details ===")
	for _, vLog := range receipt.Logs {
		if vLog.Address == contractAddress && len(vLog.Topics) > 0 {
			// Check if this is our DataSentToTarget event
			eventSignature := vLog.Topics[0]
			expectedSignature := crypto.Keccak256Hash([]byte("DataSentToTarget(address,address,bytes32,bytes32,string)"))

			if eventSignature == expectedSignature {
				fmt.Println("📡 Found DataSentToTarget event:")

				// Parse indexed parameters from topics
				if len(vLog.Topics) >= 3 {
					from := common.HexToAddress(vLog.Topics[1].Hex())
					to := common.HexToAddress(vLog.Topics[2].Hex())
					fmt.Printf("  📤 From: %s\n", from.Hex())
					fmt.Printf("  📥 To: %s\n", to.Hex())
				}

				// Parse non-indexed parameters from data
				if len(vLog.Data) > 0 {
					var event struct {
						Owner  [32]byte
						Actref [32]byte
						Topic  string
					}
					err := parsedABI.UnpackIntoInterface(&event, "DataSentToTarget", vLog.Data)
					if err != nil {
						fmt.Printf("  ❌ Error parsing event data: %v\n", err)
					} else {
						fmt.Printf("  👤 Owner: 0x%x\n", event.Owner)
						fmt.Printf("  📝 Action Ref: 0x%x\n", event.Actref)
						fmt.Printf("  🏷️  Topic: %s\n", event.Topic)
					}
				}
			}
		}
	}

	fmt.Println("\n=== Event Filtering Example ===")

	// Filter recent events
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Printf("Failed to get current block: %v", err)
	} else {
		fromBlock := currentBlock - 100 // Last 100 blocks
		if fromBlock > currentBlock {   // Handle underflow
			fromBlock = 0
		}

		fmt.Printf("🔍 Filtering events from block %d to %d...\n", fromBlock, currentBlock)

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
			fmt.Printf("📊 Found %d events:\n", len(logs))
			for i, vLog := range logs {
				fmt.Printf("  🔸 Event %d - Block: %d, TxHash: %s\n",
					i+1, vLog.BlockNumber, vLog.TxHash.Hex())
			}
		}
	}

	fmt.Println("\n🎉 Sepolia interaction completed successfully!")
	fmt.Printf("🔗 Contract on Etherscan: https://sepolia.etherscan.io/address/%s\n", contractAddress.Hex())
}

// Helper function to convert string to bytes32
func stringToBytes32(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}
