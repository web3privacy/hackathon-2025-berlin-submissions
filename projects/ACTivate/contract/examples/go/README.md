# AdminContract Go Integration Example

This Go program demonstrates how to interact with the AdminContract smart contract using Go and the go-ethereum library.

## Features

- Connect to Ethereum network (local or remote)
- Call contract functions (sendDataToTarget)
- Query contract state (getAdmin, owner)
- Parse and display events
- Filter historical events
- Handle transaction receipts and confirmations

## Prerequisites

- Go 1.21 or later
- Access to an Ethereum node (local Hardhat network or testnet)
- Contract deployed and contract address
- Private key with admin privileges on the contract

## Setup

1. **Install dependencies:**
   ```bash
   cd examples/go
   go mod tidy
   ```

2. **Set environment variables:**
   ```bash
   export RPC_URL="http://localhost:8545"  # Your Ethereum RPC endpoint
   export CONTRACT_ADDRESS="0x..."         # Deployed contract address
   export PRIVATE_KEY="0x..."              # Private key of admin account
   ```

3. **Run the example:**
   ```bash
   go run main.go
   ```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RPC_URL` | Ethereum RPC endpoint | `http://localhost:8545` |
| `CONTRACT_ADDRESS` | Deployed AdminContract address | Required |
| `PRIVATE_KEY` | Private key of admin account | Required |

## Example Output

```
=== AdminContract Go Interaction Example ===

Connected to: http://localhost:8545
Using account: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Contract address: 0x5FbDB2315678afecb367f032d93F642f64180aa3

=== Contract Information ===
Admin address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Is caller admin? true

=== Sending Data to Target ===
Target address: 0x742d35Cc6634C0532925a3b8d6968e2a4aF9c11B
Owner param: 0x4f574e45525f474f5f303031000000000000000000000000000000000000000000
Action ref: 0x4143545245465f474f5f313233000000000000000000000000000000000000000
Topic: Go Example - Blockchain Integration

Sending transaction...
Transaction hash: 0x...
Waiting for confirmation...
Transaction confirmed in block: 2
Gas used: 52341
✓ Transaction successful!

=== Event Details ===
Found DataSentToTarget event:
  From: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
  To: 0x742d35Cc6634C0532925a3b8d6968e2a4aF9c11B
  Owner: 0x4f574e45525f474f5f303031000000000000000000000000000000000000000000
  Action Ref: 0x4143545245465f474f5f313233000000000000000000000000000000000000000
  Topic: Go Example - Blockchain Integration

=== Reading Contract State ===
Contract owner: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

=== Event Filtering Example ===
Filtering events from block 0 to 2...
Found 1 events:
  Event 1 - Block: 2, TxHash: 0x...

=== Go Integration Complete ===
Successfully demonstrated AdminContract interaction from Go!
```

## Code Structure

The example demonstrates:

1. **Connection Setup**: Connect to Ethereum node and set up authentication
2. **Contract Interaction**: Create contract instance and call functions
3. **Transaction Handling**: Send transactions and wait for confirmations
4. **Event Parsing**: Parse and display contract events
5. **State Queries**: Read contract state and admin information
6. **Event Filtering**: Query historical events from the blockchain

## Key Functions

- `sendDataToTarget(target, ownerParam, actref, topic)` - Send data to target address (admin only)
- `getAdmin()` - Get the contract admin address
- `owner()` - Get the contract owner address (same as admin)

## Error Handling

The example includes proper error handling for:
- Connection failures
- Invalid private keys
- Transaction failures
- Access control violations
- Event parsing errors

## Security Notes

- Keep your private key secure and never commit it to version control
- Use environment variables or secure key management systems
- The private key account must be the contract admin to call `sendDataToTarget`
- Validate all inputs before sending transactions

## Dependencies

- `github.com/ethereum/go-ethereum` - Official Go Ethereum library
- Standard Go libraries for crypto, networking, and formatting

## Troubleshooting

1. **Connection refused**: Make sure your Ethereum node is running
2. **Invalid private key**: Ensure the private key is in hexadecimal format
3. **Transaction failed**: Check if the account has sufficient gas and is the admin
4. **Contract not found**: Verify the contract address is correct and deployed
