import { ethers } from "hardhat";

async function main() {
  console.log("Deploying Admin Contract...");

  // Get the deployer account
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with account:", deployer.address);
  console.log("Account balance:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)));

  // Deploy Admin Contract
  const AdminContract = await ethers.getContractFactory("AdminContract");
  const adminContract = await AdminContract.deploy();

  await adminContract.waitForDeployment();
  const contractAddress = await adminContract.getAddress();

  console.log("Admin Contract deployed to:", contractAddress);
  console.log("Admin address:", deployer.address);

  // Verify the deployment
  const admin = await adminContract.getAdmin();

  console.log("\n--- Deployment Verification ---");
  console.log("Contract admin:", admin);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
