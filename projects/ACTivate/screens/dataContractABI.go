package screens

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const DataContractABI = `[
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

func ParseContractABI() (abi.ABI, error) {
	parsedABI, err := abi.JSON(strings.NewReader(DataContractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	return parsedABI, nil
}
