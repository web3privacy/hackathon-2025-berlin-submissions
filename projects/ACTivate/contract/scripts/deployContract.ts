import { ethers } from "hardhat";

async function main() {
  console.log("Deploying Data Contract...");

  // Get the deployer account
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with account:", deployer.address);
  console.log("Account balance:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)));

  // Deploy Data Contract
  const DataContract = await ethers.getContractFactory("DataContract");
  const dataContract = await DataContract.deploy();

  await dataContract.waitForDeployment();
  const contractAddress = await dataContract.getAddress();

  console.log("Data Contract deployed to:", contractAddress);
  console.log("Deployer address:", deployer.address);

  // Verify the deployment
  console.log("\n--- Deployment Verification ---");
  console.log("Contract deployed successfully!");
  console.log("Contract is public - anyone can call sendDataToTarget function");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
