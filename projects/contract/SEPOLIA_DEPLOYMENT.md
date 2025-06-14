# AdminContract Sepolia Deployment Guide

This guide explains how to deploy the AdminContract to the Sepolia testnet and use it with Go.

## 📋 Prerequisites

1. **Node.js & npm** - For running the deployment script
2. **Sepolia ETH** - Get free testnet ETH from [Sepolia Faucet](https://sepoliafaucet.com/)
3. **RPC Provider** - Infura, Alchemy, or another Ethereum RPC provider
4. **Private Key** - From MetaMask or another wallet (keep it secure!)
5. **Go 1.21+** - For running Go integration examples

## 🔧 Setup

### 1. Environment Configuration

Copy the environment template:
```bash
cp .env.example .env
```

Edit `.env` with your values:
```bash
# Get from https://infura.io/ or https://dashboard.alchemy.com/
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID

# Your wallet private key (keep secure!)
PRIVATE_KEY=0x1234567890abcdef...

# Optional: For contract verification
ETHERSCAN_API_KEY=YOUR_ETHERSCAN_API_KEY
```

### 2. Get Testnet ETH

Visit [Sepolia Faucet](https://sepoliafaucet.com/) and request testnet ETH for your wallet address.

### 3. Install Dependencies

```bash
npm install
```

## 🚀 Deployment

### Deploy to Sepolia

```bash
npm run deploy:sepolia
```

This will:
- ✅ Deploy the AdminContract to Sepolia testnet
- ✅ Save deployment info in multiple formats
- ✅ Generate Go-compatible configuration files
- ✅ Provide Etherscan links for verification

### Expected Output

```
=== AdminContract Sepolia Deployment ===

📡 Network: sepolia (Chain ID: 11155111)
🔗 RPC URL: https://sepolia.infura.io/v3/...
👤 Deployer address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
💰 Deployer balance: 0.5 ETH
⛽ Gas Price: 2.5 gwei

🚀 Deploying AdminContract...
📊 Estimated gas: 567,890
📋 Transaction submitted: 0x123...
⏳ Waiting for deployment confirmation...

✅ Deployment successful!
📍 Contract address: 0x1234567890abcdef1234567890abcdef12345678
🔗 Transaction hash: 0x123...
📦 Block number: 4567890
⛽ Gas used: 543,210

🔍 Verifying contract deployment...
✅ Contract admin: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
✅ Admin verification: PASSED

💾 Deployment info saved to: deployments/sepolia-deployment.json
🐹 Go configuration saved to: deployments/sepolia-config.go
📄 Environment file saved to: deployments/sepolia.env
🔧 Go constants file saved to: deployments/constants.go

🎉 Deployment completed successfully!
```

## 📁 Generated Files

After deployment, you'll find these files in the `deployments/` directory:

### 1. `sepolia-deployment.json`
Complete deployment information in JSON format.

### 2. `sepolia-config.go`
Go struct with all deployment data embedded.

### 3. `sepolia.env`
Environment variables for easy loading.

### 4. `constants.go`
Go constants file ready to copy to your project.

## 🔍 Contract Verification (Optional)

Verify your contract on Etherscan:

```bash
npx hardhat verify --network sepolia CONTRACT_ADDRESS
```

Or using the npm script:
```bash
npm run verify:sepolia CONTRACT_ADDRESS
```

## 🐹 Go Integration

### Method 1: Using Environment Variables

Load the generated environment file:
```bash
source deployments/sepolia.env
cd examples/go
go run sepolia-interaction.go
```

### Method 2: Using Generated Constants

Copy the constants to your Go project:
```bash
cp deployments/constants.go /path/to/your/go/project/
```

Update the constants in your code and run:
```go
package main

import (
    // ... your imports
)

// Use the generated constants
func main() {
    rpcURL := SepoliaRPCURL
    contractAddress := AdminContractAddress
    privateKey := PrivateKey
    
    // ... rest of your code
}
```

### Method 3: Load Configuration Dynamically

Use the generated Go configuration:
```go
package main

import (
    "path/to/your/deployments"
)

func main() {
    config := deployments.GetSepoliaConfig()
    
    client, _ := ethclient.Dial(config.RPCUrl)
    contract := common.HexToAddress(config.ContractAddress)
    
    // ... rest of your code
}
```

## 🧪 Testing the Deployment

### Test with TypeScript

```bash
CONTRACT_ADDRESS=0x... npm run interact
```

### Test with Go

```bash
cd examples/go
export RPC_URL="https://sepolia.infura.io/v3/YOUR_PROJECT_ID"
export CONTRACT_ADDRESS="0x..."
export PRIVATE_KEY="0x..."
go run sepolia-interaction.go
```

## 📊 Monitoring

### View on Etherscan

- Contract: `https://sepolia.etherscan.io/address/CONTRACT_ADDRESS`
- Deployment Transaction: `https://sepolia.etherscan.io/tx/TX_HASH`

### Check Contract State

```bash
# Using cast (if you have Foundry installed)
cast call CONTRACT_ADDRESS "owner()" --rpc-url $SEPOLIA_RPC_URL

# Or check on Etherscan directly
```

## 🔐 Security Best Practices

### Environment Variables
- ✅ Never commit `.env` files to version control
- ✅ Use different private keys for testnet and mainnet
- ✅ Consider using hardware wallets for mainnet deployments
- ✅ Rotate keys periodically

### Private Key Management
```bash
# Good: Use environment variables
export PRIVATE_KEY="0x..."

# Better: Use secure key management systems
# - AWS KMS
# - HashiCorp Vault  
# - Hardware Security Modules (HSM)
```

### Code Security
- ✅ Always verify contract source code on Etherscan
- ✅ Test thoroughly on testnet before mainnet deployment
- ✅ Use multi-signature wallets for production contracts
- ✅ Implement proper access controls

## 🐛 Troubleshooting

### Common Issues

1. **Insufficient Balance**
   ```
   Error: insufficient funds for gas * price + value
   ```
   **Solution**: Get more Sepolia ETH from the faucet

2. **Invalid Private Key**
   ```
   Error: invalid private key
   ```
   **Solution**: Ensure private key is valid hex (with or without 0x prefix)

3. **Network Connection**
   ```
   Error: network connection failed
   ```
   **Solution**: Check your RPC URL and internet connection

4. **Gas Estimation Failed**
   ```
   Error: cannot estimate gas
   ```
   **Solution**: Check contract parameters and account permissions

### Debug Commands

```bash
# Check network connection
curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' $SEPOLIA_RPC_URL

# Check account balance
cast balance YOUR_ADDRESS --rpc-url $SEPOLIA_RPC_URL

# Check contract code
cast code CONTRACT_ADDRESS --rpc-url $SEPOLIA_RPC_URL
```

## 📚 Additional Resources

- [Sepolia Testnet Info](https://sepolia.dev/)
- [Sepolia Faucet](https://sepoliafaucet.com/)
- [Etherscan Sepolia](https://sepolia.etherscan.io/)
- [Infura Documentation](https://docs.infura.io/)
- [Alchemy Documentation](https://docs.alchemy.com/)
- [go-ethereum Documentation](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings)

## 🎯 Next Steps

After successful deployment:

1. **✅ Verify Contract** - Submit source code to Etherscan
2. **✅ Test Interactions** - Use the provided scripts to test functionality
3. **✅ Build Frontend** - Create a web interface for your contract
4. **✅ Set up Monitoring** - Track events and contract usage
5. **✅ Plan Mainnet** - Prepare for production deployment

---

**🚨 Remember: Never use testnet private keys on mainnet!**
