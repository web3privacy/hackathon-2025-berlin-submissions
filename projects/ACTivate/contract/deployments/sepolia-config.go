package main

import (
	"encoding/json"
	"log"
)

// AdminContractConfig contains all deployment information for Sepolia testnet
type AdminContractConfig struct {
	Network         string      `json:"network"`
	ChainID         int64       `json:"chainId"`
	ContractAddress string      `json:"contractAddress"`
	ContractABI     interface{} `json:"contractABI"`
	DeployerAddress string      `json:"deployerAddress"`
	PrivateKey      string      `json:"privateKey"`
	RPCUrl          string      `json:"rpcUrl"`
	BlockNumber     int64       `json:"blockNumber"`
	TransactionHash string      `json:"transactionHash"`
	GasUsed         string      `json:"gasUsed"`
	DeployedAt      string      `json:"deployedAt"`
}

// GetSepoliaConfig returns the deployment configuration for Sepolia testnet
func GetSepoliaConfig() *AdminContractConfig {
	configJSON := `{
  "network": "sepolia",
  "chainId": 11155111,
  "contractAddress": "0x442f8f596045BcB87E3B38C58A42F40797F81F7E",
  "contractABI": [
    {
      "inputs": [],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnableInvalidOwner",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "OwnableUnauthorizedAccount",
      "type": "error"
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
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
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
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
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
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ],
  "deployerAddress": "0x5225c07Ec3ba1D5fE360459fE5B9C2Db28b35c9B",
  "privateKey": "0xdf0208788c196f1d440532cebf35c4beda88df5dccc7f1f15492cd5ea118e56c",
  "rpcUrl": "https://eth-sepolia.g.alchemy.com/v2/atcICv4EFi9hXKew1D4LvnH36cm5-96S",
  "blockNumber": 0,
  "transactionHash": "0x5bb527e6528bd19d108f2b3a4bbc86a0aa029c9707842927ae8e7df9fa7ed35b",
  "gasUsed": "288805",
  "deployedAt": "2025-06-14T11:37:49.207Z"
}`

	var config AdminContractConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &config
}
