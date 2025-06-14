import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import * as dotenv from "dotenv";

// Load environment variables
dotenv.config();

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.20",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  networks: {
    hardhat: {
      chainId: 1337,
    },
    localhost: {
      url: "http://127.0.0.1:8545",
    },
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL || "https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      chainId: 11155111,
      gasPrice: 20000000000, // 20 gwei
    },
    gnosis: {
      url: process.env.GNOSIS_RPC_URL || "https://rpc.gnosischain.com",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      chainId: 100,
      gasPrice: 1000000000, // 1 gwei (cheaper on Gnosis)
    },
    chiado: {
      url: process.env.CHIADO_RPC_URL || "https://rpc.chiadochain.net",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      chainId: 10200,
      gasPrice: 1000000000, // 1 gwei
    },
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
};

export default config;