import { ethers } from "hardhat";

async function main() {
  console.log("=== Zero Address Validation Test ===\n");

  // Deploy a fresh contract
  const DataContract = await ethers.getContractFactory("DataContract");
  const dataContract = await DataContract.deploy();
  await dataContract.waitForDeployment();
  
  const contractAddress = await dataContract.getAddress();
  console.log("Fresh contract deployed at:", contractAddress);

  // Get a test account
  const [, , , user3] = await ethers.getSigners();
  
  // Test zero address validation immediately
  console.log("\nTesting zero address validation...");
  const zeroAddress = "0x0000000000000000000000000000000000000000";
  console.log("Zero address:", zeroAddress);
  
  try {
    const tx = await dataContract.connect(user3).sendDataToTarget(
      zeroAddress,
      ethers.keccak256(ethers.toUtf8Bytes("ERROR_TEST")),
      ethers.keccak256(ethers.toUtf8Bytes("ERROR_ACTION")),
      "This should fail"
    );
    
    console.log("Transaction hash:", tx.hash);
    const receipt = await tx.wait();
    console.log("Transaction receipt status:", receipt?.status);
    
    if (receipt?.status === 0) {
      console.log("✅ Transaction failed as expected (receipt status: 0)");
    } else {
      console.log("❌ Transaction succeeded but should have failed!");
    }
    
  } catch (error) {
    console.log("✅ Error caught as expected:");
    const errorMessage = error instanceof Error ? error.message : String(error);
    console.log("   Error:", errorMessage.split('\n')[0]);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("Script failed:", error);
    process.exit(1);
  });
