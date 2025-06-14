const { ethers } = require('ethers');
require('dotenv').config();

async function checkBalance() {
  // Connect to Sepolia network
  const provider = new ethers.JsonRpcProvider(process.env.SEPOLIA_RPC_URL);
  
  // Get wallet from private key
  const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider);
  
  console.log('🔍 Checking Sepolia ETH Balance');
  console.log('================================');
  console.log(`Address: ${wallet.address}`);
  
  try {
    const balance = await provider.getBalance(wallet.address);
    const balanceInEth = ethers.formatEther(balance);
    
    console.log(`Balance: ${balanceInEth} ETH`);
    
    if (parseFloat(balanceInEth) > 0) {
      console.log('✅ You have ETH! Ready to deploy.');
    } else {
      console.log('❌ No ETH found. Please fund your address from a faucet:');
      console.log('- https://sepoliafaucet.com/');
      console.log('- https://faucet.sepolia.dev/');
      console.log('- https://www.alchemy.com/faucets/ethereum-sepolia');
    }
  } catch (error) {
    console.error('Error checking balance:', error.message);
  }
}

checkBalance().catch(console.error);
