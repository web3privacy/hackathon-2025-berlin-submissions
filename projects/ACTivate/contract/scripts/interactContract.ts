import { ethers } from "hardhat";
import { DataContract } from "../typechain-types";

async function main() {
  console.log("=== DataContract Interaction Script ===\n");

  // Get the deployed contract address
  const contractAddress = process.env.CONTRACT_ADDRESS;
  if (!contractAddress) {
    console.error("Error: CONTRACT_ADDRESS environment variable not set");
    console.log("Please set CONTRACT_ADDRESS to the deployed contract address");
    console.log("Example: CONTRACT_ADDRESS=0x... npm run interact");
    process.exit(1);
  }

  // Get signers
  const [deployer, user1] = await ethers.getSigners();
  console.log("Deployer address:", deployer.address);
  console.log("User1 address:", user1.address);

  // Connect to the deployed contract
  const DataContract = await ethers.getContractFactory("DataContract");
  const dataContract = DataContract.attach(contractAddress) as DataContract;

  console.log("\n=== Contract Information ===");
  console.log("Contract address:", contractAddress);
  console.log("Public access enabled - anyone can call functions");
  console.log("Current caller:", deployer.address);

  console.log("\n=== Sending Data to Target ===");
  
  // Example data to send
  const targetAddress = user1.address;
  const ownerParam = ethers.encodeBytes32String("OWNER_001");
  const actref = ethers.encodeBytes32String("ACTION_REF_123");
  const topic = "Test Topic - Contract Interaction";

  console.log("Target address:", targetAddress);
  console.log("Owner param:", ownerParam);
  console.log("Action reference:", actref);
  console.log("Topic:", topic);

  try {
    // Send data to target (anyone can do this now)
    console.log("\nSending data to target...");
    const tx = await dataContract.sendDataToTarget(
      targetAddress,
      ownerParam,
      actref,
      topic
    );
    
    console.log("Transaction hash:", tx.hash);
    console.log("Waiting for confirmation...");
    
    const receipt = await tx.wait();
    if (!receipt) {
      console.error("Transaction receipt is null");
      return;
    }
    
    console.log("Transaction confirmed in block:", receipt.blockNumber);

    // Check the event logs
    console.log(`Found ${receipt.logs.length} logs in transaction`);
    for (let i = 0; i < receipt.logs.length; i++) {
      const log = receipt.logs[i];
      try {
        const parsedLog = dataContract.interface.parseLog({
          topics: log.topics,
          data: log.data
        });
        if (parsedLog && parsedLog.name === 'DataSentToTarget') {
          console.log("\n=== Event Details ===");
          console.log("From:", parsedLog.args[0]);
          console.log("To:", parsedLog.args[1]);
          console.log("Owner:", parsedLog.args[2]);
          console.log("Action Ref:", parsedLog.args[3]);
          console.log("Topic:", parsedLog.args[4]);
        }
      } catch (error) {
        // Log might not be from our contract, ignore
      }
    }

  } catch (error: any) {
    console.error("Error sending data:", error.message);
  }

  console.log("\n=== Testing Multi-User Access ===");
  
  try {
    // Test calling from different account (should now work)
    console.log("Testing call from different account...");
    const userContract = dataContract.connect(user1);
    
    const tx = await userContract.sendDataToTarget(
      deployer.address,
      ethers.encodeBytes32String("USER1_CALL"),
      ethers.encodeBytes32String("PUBLIC_ACCESS"),
      "Call from user1 - should work now"
    );
    
    const receipt = await tx.wait();
    if (receipt) {
      console.log("✓ Multi-user access working correctly!");
      console.log("Transaction confirmed in block:", receipt.blockNumber);
    }
    
  } catch (error: any) {
    console.log("✗ Multi-user access failed:", error.message);
  }

  console.log("\n=== Multiple Data Sends ===");
  
  // Send multiple data entries
  const dataEntries = [
    {
      target: user1.address,
      owner: ethers.encodeBytes32String("OWNER_002"),
      actref: ethers.encodeBytes32String("REF_BATCH_001"),
      topic: "Batch Entry 1"
    },
    {
      target: deployer.address,
      owner: ethers.encodeBytes32String("OWNER_003"),
      actref: ethers.encodeBytes32String("REF_BATCH_002"),
      topic: "Batch Entry 2"
    },
    {
      target: user1.address,
      owner: ethers.encodeBytes32String("OWNER_004"),
      actref: ethers.encodeBytes32String("REF_BATCH_003"),
      topic: "Final Batch Entry"
    }
  ];

  for (let i = 0; i < dataEntries.length; i++) {
    const entry = dataEntries[i];
    console.log(`\nSending batch entry ${i + 1}/${dataEntries.length}...`);
    
    try {
      const tx = await dataContract.sendDataToTarget(
        entry.target,
        entry.owner,
        entry.actref,
        entry.topic
      );
      
      const receipt = await tx.wait();
      if (receipt) {
        console.log(`✓ Batch entry ${i + 1} sent successfully (Block: ${receipt.blockNumber})`);
      } else {
        console.log(`✗ Batch entry ${i + 1} failed: Receipt is null`);
      }
      
    } catch (error: any) {
      console.error(`✗ Batch entry ${i + 1} failed:`, error.message);
    }
  }

  console.log("\n=== Interaction Complete ===");
  console.log("Contract interaction script finished successfully!");
}

// Execute the script
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("Script failed:", error);
    process.exit(1);
  });