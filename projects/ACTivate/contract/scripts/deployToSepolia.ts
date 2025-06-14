import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

// Configuration interface for Go integration
interface DeploymentConfig {
  network: string;
  chainId: number;
  contractAddress: string;
  contractABI: any[]; // eslint-disable-line @typescript-eslint/no-explicit-any
  deployerAddress: string;
  privateKey: string;
  rpcUrl: string;
  blockNumber: number;
  transactionHash: string;
  gasUsed: string;
  deployedAt: string;
}

async function main(): Promise<void> {
  console.log("=== AdminContract Sepolia Deployment ===\n");

  // Check environment variables
  const privateKey = process.env.PRIVATE_KEY;
  const rpcUrl = process.env.SEPOLIA_RPC_URL;

  if (!privateKey) {
    console.error("❌ Error: PRIVATE_KEY environment variable not set");
    console.log("Please set your private key:");
    console.log("export PRIVATE_KEY=0x...");
    process.exit(1);
  }

  if (!rpcUrl) {
    console.error("❌ Error: SEPOLIA_RPC_URL environment variable not set");
    console.log("Please set your Sepolia RPC URL:");
    console.log("export SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID");
    console.log("or");
    console.log("export SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY");
    process.exit(1);
  }

  // Get network info
  const network = await ethers.provider.getNetwork();
  console.log(`📡 Network: ${network.name} (Chain ID: ${network.chainId})`);
  console.log(`🔗 RPC URL: ${rpcUrl}`);

  // Get deployer account
  const [deployer] = await ethers.getSigners();
  console.log(`👤 Deployer address: ${deployer.address}`);

  // Check balance
  const balance = await ethers.provider.getBalance(deployer.address);
  const balanceInEth = ethers.formatEther(balance);
  console.log(`💰 Deployer balance: ${balanceInEth} ETH`);

  if (parseFloat(balanceInEth) < 0.01) {
    console.warn("⚠️  Warning: Low balance! You might need more ETH for deployment.");
    console.log("💡 Get Sepolia ETH from: https://sepoliafaucet.com/");
  }

  // Get current gas price
  const feeData = await ethers.provider.getFeeData();
  console.log(`⛽ Gas Price: ${ethers.formatUnits(feeData.gasPrice || 0, "gwei")} gwei`);

  console.log("\n🚀 Deploying AdminContract...");

  // Deploy the contract
  const AdminContract = await ethers.getContractFactory("AdminContract");
  
  // Estimate gas for deployment
  const deploymentData = await AdminContract.getDeployTransaction();
  const estimatedGas = await ethers.provider.estimateGas(deploymentData);
  console.log(`📊 Estimated gas: ${estimatedGas.toString()}`);

  const adminContract = await AdminContract.deploy();
  console.log(`📋 Transaction submitted: ${adminContract.deploymentTransaction()?.hash}`);
  
  // Wait for deployment
  console.log("⏳ Waiting for deployment confirmation...");
  await adminContract.waitForDeployment();

  const contractAddress = await adminContract.getAddress();
  const deploymentReceipt = adminContract.deploymentTransaction();
  
  console.log("\n✅ Deployment successful!");
  console.log(`📍 Contract address: ${contractAddress}`);
  console.log(`🔗 Transaction hash: ${deploymentReceipt?.hash}`);
  console.log(`📦 Block number: ${deploymentReceipt?.blockNumber}`);
  console.log(`⛽ Gas used: ${deploymentReceipt?.gasLimit?.toString()}`);

  // Verify contract functions
  console.log("\n🔍 Verifying contract deployment...");
  try {
    const adminAddress = await adminContract.owner();
    console.log(`✅ Contract admin: ${adminAddress}`);
    console.log(`✅ Admin verification: ${adminAddress === deployer.address ? "PASSED" : "FAILED"}`);
  } catch (error) {
    console.error("❌ Contract verification failed:", error instanceof Error ? error.message : String(error));
  }

  // Load contract ABI from artifacts
  const artifactPath = path.join(__dirname, "..", "artifacts", "contracts", "AdminContract.sol", "AdminContract.json");
  
  let artifact: { abi: any[] }; // eslint-disable-line @typescript-eslint/no-explicit-any
  try {
    artifact = JSON.parse(fs.readFileSync(artifactPath, "utf8"));
  } catch (error) {
    console.error("❌ Failed to load contract artifact:", error instanceof Error ? error.message : String(error));
    process.exit(1);
  }

  // Create deployment configuration for Go
  const deploymentConfig: DeploymentConfig = {
    network: network.name,
    chainId: Number(network.chainId),
    contractAddress: contractAddress,
    contractABI: artifact.abi,
    deployerAddress: deployer.address,
    privateKey: privateKey,
    rpcUrl: rpcUrl,
    blockNumber: deploymentReceipt?.blockNumber || 0,
    transactionHash: deploymentReceipt?.hash || "",
    gasUsed: deploymentReceipt?.gasLimit?.toString() || "0",
    deployedAt: new Date().toISOString(),
  };

  // Save deployment info in multiple formats
  const deploymentsDir = path.join(__dirname, "..", "deployments");
  try {
    if (!fs.existsSync(deploymentsDir)) {
      fs.mkdirSync(deploymentsDir, { recursive: true });
    }

    // 1. JSON format for general use
    const jsonPath = path.join(deploymentsDir, `sepolia-deployment.json`);
    fs.writeFileSync(jsonPath, JSON.stringify(deploymentConfig, null, 2));
    console.log(`💾 Deployment info saved to: ${jsonPath}`);

    // 2. Go configuration file
    const goConfigPath = path.join(deploymentsDir, `sepolia-config.go`);
    const goConfig = generateGoConfig(deploymentConfig);
    fs.writeFileSync(goConfigPath, goConfig);
    console.log(`🐹 Go configuration saved to: ${goConfigPath}`);

    // 3. Environment file for easy loading
    const envPath = path.join(deploymentsDir, `sepolia.env`);
    const envContent = generateEnvFile(deploymentConfig);
    fs.writeFileSync(envPath, envContent);
    console.log(`📄 Environment file saved to: ${envPath}`);

    // 4. Go constants file
    const goConstantsPath = path.join(deploymentsDir, `constants.go`);
    const goConstants = generateGoConstants(deploymentConfig);
    fs.writeFileSync(goConstantsPath, goConstants);
    console.log(`🔧 Go constants file saved to: ${goConstantsPath}`);
  } catch (error) {
    console.error("❌ Failed to save deployment files:", error instanceof Error ? error.message : String(error));
    console.log("⚠️  Deployment was successful but files could not be saved.");
  }

  console.log("\n🎉 Deployment completed successfully!");
  console.log("\n📋 Next steps:");
  console.log("1. Verify contract on Etherscan (optional):");
  console.log(`   npx hardhat verify --network sepolia ${contractAddress}`);
  console.log("\n2. Test contract interaction:");
  console.log(`   CONTRACT_ADDRESS=${contractAddress} npm run interact`);
  console.log("\n3. Use in Go code:");
  console.log(`   Copy deployments/constants.go to your Go project`);
  console.log(`   Or source deployments/sepolia.env in your environment`);

  console.log(`\n🔗 View on Etherscan: https://sepolia.etherscan.io/address/${contractAddress}`);
}

function generateGoConfig(config: DeploymentConfig): string {
  return `package main

import (
	"encoding/json"
	"log"
)

// AdminContractConfig contains all deployment information for Sepolia testnet
type AdminContractConfig struct {
	Network         string      \`json:"network"\`
	ChainID         int64       \`json:"chainId"\`
	ContractAddress string      \`json:"contractAddress"\`
	ContractABI     interface{} \`json:"contractABI"\`
	DeployerAddress string      \`json:"deployerAddress"\`
	PrivateKey      string      \`json:"privateKey"\`
	RPCUrl          string      \`json:"rpcUrl"\`
	BlockNumber     int64       \`json:"blockNumber"\`
	TransactionHash string      \`json:"transactionHash"\`
	GasUsed         string      \`json:"gasUsed"\`
	DeployedAt      string      \`json:"deployedAt"\`
}

// GetSepoliaConfig returns the deployment configuration for Sepolia testnet
func GetSepoliaConfig() *AdminContractConfig {
	configJSON := \`${JSON.stringify(config, null, 2)}\`
	
	var config AdminContractConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
	
	return &config
}

// AdminContractABI contains the contract ABI as a JSON string
const AdminContractABI = \`${JSON.stringify(config.contractABI)}\`
`;
}

function generateEnvFile(config: DeploymentConfig): string {
  return `# AdminContract Sepolia Deployment Configuration
# Generated on ${config.deployedAt}

# Network Configuration
NETWORK=${config.network}
CHAIN_ID=${config.chainId}
RPC_URL=${config.rpcUrl}

# Contract Information
CONTRACT_ADDRESS=${config.contractAddress}
DEPLOYER_ADDRESS=${config.deployerAddress}

# Transaction Information
DEPLOYMENT_TX_HASH=${config.transactionHash}
DEPLOYMENT_BLOCK=${config.blockNumber}
GAS_USED=${config.gasUsed}

# Private Key (Keep secure!)
PRIVATE_KEY=${config.privateKey}

# Etherscan Links
ETHERSCAN_CONTRACT_URL=https://sepolia.etherscan.io/address/${config.contractAddress}
ETHERSCAN_TX_URL=https://sepolia.etherscan.io/tx/${config.transactionHash}
`;
}

function generateGoConstants(config: DeploymentConfig): string {
  return `package main

// AdminContract Sepolia deployment constants
// Generated on ${config.deployedAt}

const (
	// Network Configuration
	SepoliaChainID = ${config.chainId}
	SepoliaRPCURL  = "${config.rpcUrl}"
	
	// Contract Information
	AdminContractAddress = "${config.contractAddress}"
	DeployerAddress     = "${config.deployerAddress}"
	
	// Transaction Information
	DeploymentTxHash  = "${config.transactionHash}"
	DeploymentBlock   = ${config.blockNumber}
	GasUsed          = "${config.gasUsed}"
)

// Private key - Keep this secure and consider using environment variables
const PrivateKey = "${config.privateKey}"

// Contract ABI as JSON string
const AdminContractABI = \`${JSON.stringify(config.contractABI)}\`

// Etherscan URLs
const (
	EtherscanContractURL = "https://sepolia.etherscan.io/address/${config.contractAddress}"
	EtherscanTxURL      = "https://sepolia.etherscan.io/tx/${config.transactionHash}"
)
`;
}

// Execute the deployment
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Deployment failed:", error);
    process.exit(1);
  });
