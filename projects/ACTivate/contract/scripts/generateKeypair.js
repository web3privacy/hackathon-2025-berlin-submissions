const { ethers } = require('ethers');

// Generate a random wallet
const wallet = ethers.Wallet.createRandom();

console.log('Generated Ethereum Keypair:');
console.log('Private Key:', wallet.privateKey);
console.log('Address:', wallet.address);
console.log('');
console.log('IMPORTANT: Save this private key securely!');
console.log('You will need Sepolia ETH at this address to deploy the contract.');
console.log('');
console.log('Get Sepolia ETH from faucets:');
console.log('- https://sepoliafaucet.com/');
console.log('- https://faucet.sepolia.dev/');
console.log('- https://www.alchemy.com/faucets/ethereum-sepolia');
