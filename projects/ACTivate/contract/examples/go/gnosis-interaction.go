package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// Configuration constants for Gnosis Chain
const (
	GnosisChainID     = 100
	GnosisRPCURL      = "https://rpc.gnosischain.com"
	GnosisExplorerURL = "https://gnosisscan.io"

	ChiadoChainID     = 10200
	ChiadoRPCURL      = "https://rpc.chiadochain.net"
	ChiadoExplorerURL = "https://gnosis-chiado.blockscout.com"
)

func main() {
	// Load environment variables
	err := godotenv.Load("deployments/gnosis.env")
	if err != nil {
		log.Println("Warning: Could not load gnosis.env file, trying .env")
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Configuration
	rpcURL := getEnvOrDefault("GNOSIS_RPC_URL", GnosisRPCURL)
	contractAddressStr := os.Getenv("GNOSIS_CONTRACT_ADDRESS")
	privateKeyStr := os.Getenv("PRIVATE_KEY")

	if contractAddressStr == "" {
		log.Fatal("GNOSIS_CONTRACT_ADDRESS not set. Please deploy the contract first with: npm run deploy:gnosis")
	}

	if privateKeyStr == "" {
		log.Fatal("PRIVATE_KEY not set")
	}

	// Remove '0x' prefix if present
	if len(privateKeyStr) > 2 && privateKeyStr[:2] == "0x" {
		privateKeyStr = privateKeyStr[2:]
	}

	fmt.Println("🌐 Gnosis Chain AdminContract Interaction")
	fmt.Println("=========================================")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Contract Address: %s\n", contractAddressStr)

	// Connect to Gnosis Chain
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to Gnosis Chain:", err)
	}
	defer client.Close()

	// Verify we're connected to the right network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("Failed to get chain ID:", err)
	}

	fmt.Printf("Connected to chain ID: %s\n", chainID.String())
	if chainID.Cmp(big.NewInt(GnosisChainID)) != 0 && chainID.Cmp(big.NewInt(ChiadoChainID)) != 0 {
		log.Printf("Warning: Connected to unexpected chain ID %s", chainID.String())
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Interacting from address: %s\n", fromAddress.Hex())

	// Check balance
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal("Failed to get balance:", err)
	}
	fmt.Printf("Account balance: %s xDAI\n", formatEther(balance))

	// Parse contract address
	contractAddress := common.HexToAddress(contractAddressStr)

	// Create contract instance (simplified ABI for demonstration)
	// In a real application, you would generate Go bindings from the ABI
	contractABI := `[
		{
			"inputs": [
				{"internalType": "address", "name": "target", "type": "address"},
				{"internalType": "bytes32", "name": "ownerParam", "type": "bytes32"},
				{"internalType": "bytes32", "name": "actref", "type": "bytes32"},
				{"internalType": "string", "name": "topic", "type": "string"}
			],
			"name": "sendDataToTarget",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"anonymous": false,
			"inputs": [
				{"indexed": true, "internalType": "address", "name": "from", "type": "address"},
				{"indexed": true, "internalType": "address", "name": "to", "type": "address"},
				{"indexed": false, "internalType": "bytes32", "name": "owner", "type": "bytes32"},
				{"indexed": false, "internalType": "bytes32", "name": "actref", "type": "bytes32"},
				{"indexed": false, "internalType": "string", "name": "topic", "type": "string"}
			],
			"name": "DataSentToTarget",
			"type": "event"
		}
	]`

	// Get the current gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Failed to get gas price:", err)
	}

	// Gnosis Chain typically has lower gas prices
	fmt.Printf("Suggested gas price: %s gwei\n", formatGwei(gasPrice))

	// Get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}

	// Create transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("Failed to create transactor:", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(100000) // Adjust as needed
	auth.GasPrice = gasPrice

	fmt.Println("\n📋 Contract Interaction Example")
	fmt.Println("For actual interaction, you would:")
	fmt.Println("1. Generate Go bindings: abigen --abi DataContract.abi --pkg main --type DataContract --out DataContract.go")
	fmt.Println("2. Use the generated bindings to interact with the contract")
	fmt.Println("3. Call contract methods like sendDataToTarget()")

	explorerURL := GnosisExplorerURL
	if chainID.Cmp(big.NewInt(ChiadoChainID)) == 0 {
		explorerURL = ChiadoExplorerURL
	}

	fmt.Printf("\n🔗 Useful Links:\n")
	fmt.Printf("- Contract Explorer: %s/address/%s\n", explorerURL, contractAddress.Hex())
	fmt.Printf("- Account Explorer: %s/address/%s\n", explorerURL, fromAddress.Hex())

	if chainID.Cmp(big.NewInt(GnosisChainID)) == 0 {
		fmt.Println("- Bridge xDAI: https://bridge.gnosischain.com/")
		fmt.Println("- Gnosis Safe: https://safe.gnosis.io/")
	} else {
		fmt.Println("- Chiado Faucet: https://gnosisfaucet.com/")
	}

	fmt.Println("\n✅ Connection to Gnosis Chain successful!")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func formatEther(wei *big.Int) string {
	ether := new(big.Float)
	ether.SetString(wei.String())
	ether.Quo(ether, big.NewFloat(1e18))
	return ether.Text('f', 6)
}

func formatGwei(wei *big.Int) string {
	gwei := new(big.Float)
	gwei.SetString(wei.String())
	gwei.Quo(gwei, big.NewFloat(1e9))
	return gwei.Text('f', 2)
}
