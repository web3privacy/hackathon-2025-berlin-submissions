const { ethers } = require('ethers');
require('dotenv').config();

async function checkBalance(networkName, rpcUrl, chainId) {
  console.log(`🔍 Checking ${networkName} Balance`);
  console.log('================================');
  
  try {
    // Connect to the network
    const provider = new ethers.JsonRpcProvider(rpcUrl);
    
    // Get wallet from private key
    const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider);
    
    console.log(`Address: ${wallet.address}`);
    console.log(`Network: ${networkName} (Chain ID: ${chainId})`);
    console.log(`RPC URL: ${rpcUrl}`);
    
    const balance = await provider.getBalance(wallet.address);
    const balanceInEth = ethers.formatEther(balance);
    
    console.log(`Balance: ${balanceInEth} xDAI`);
    
    if (parseFloat(balanceInEth) > 0) {
      console.log('✅ You have xDAI! Ready to deploy.');
    } else {
      console.log('❌ No xDAI found. Please fund your address:');
      if (networkName === 'Gnosis Chain') {
        console.log('- Bridge ETH to xDAI: https://bridge.gnosischain.com/');
        console.log('- Gnosis Faucet: https://gnosisfaucet.com/');
      } else if (networkName === 'Chiado Testnet') {
        console.log('- Chiado Faucet: https://gnosisfaucet.com/');
      }
    }
    
    console.log('');
    
  } catch (error) {
    console.error(`Error checking ${networkName} balance:`, error.message);
    console.log('');
  }
}

async function main() {
  const networks = [
    {
      name: 'Gnosis Chain',
      rpcUrl: process.env.GNOSIS_RPC_URL || 'https://rpc.gnosischain.com',
      chainId: 100
    },
    {
      name: 'Chiado Testnet',
      rpcUrl: process.env.CHIADO_RPC_URL || 'https://rpc.chiadochain.net',
      chainId: 10200
    }
  ];

  console.log('🌐 Checking Gnosis Chain Networks');
  console.log('=================================\n');

  for (const network of networks) {
    await checkBalance(network.name, network.rpcUrl, network.chainId);
  }
}

main().catch(console.error);
