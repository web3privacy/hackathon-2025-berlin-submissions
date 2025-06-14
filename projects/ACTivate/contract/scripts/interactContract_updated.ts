import { ethers } from "hardhat";

async function main() {
  console.log("=== Admin Contract Interaction Demo ===\n");

  // Get signers
  const [admin, user1, user2] = await ethers.getSigners();
  console.log("Admin address:", admin.address);
  console.log("User1 address:", user1.address);
  console.log("User2 address:", user2.address);

  // Deploy the contract (in a real scenario, you'd connect to an existing deployment)
  const AdminContract = await ethers.getContractFactory("AdminContract");
  const adminContract = await AdminContract.deploy();
  await adminContract.waitForDeployment();
  
  const contractAddress = await adminContract.getAddress();
  console.log("Contract deployed at:", contractAddress);
  console.log("Contract admin:", await adminContract.getAdmin(), "\n");

  // Demonstrate sending data to target
  console.log("--- Sending Data to Target ---");
  const ownerParam = ethers.keccak256(ethers.toUtf8Bytes("transaction_owner_12345"));
  const actref = ethers.keccak256(ethers.toUtf8Bytes("action_ref_xyz"));
  const topic = "Payment Transaction";

  console.log("Sending data to:", user1.address);
  console.log("Owner param:", ownerParam);
  console.log("Action ref:", actref);
  console.log("Topic:", topic);

  // Send data to target
  const tx = await adminContract.connect(admin).sendDataToTarget(
    user1.address,
    ownerParam,
    actref,
    topic
  );

  // Wait for transaction and get receipt
  const receipt = await tx.wait();
  console.log("Transaction hash:", tx.hash);

  // Parse the event
  const iface = adminContract.interface;
  const logs = receipt?.logs.map((log: any) => {
    try {
      return iface.parseLog({ topics: log.topics, data: log.data });
    } catch {
      return null;
    }
  }).filter((log: any) => log !== null);

  if (logs && logs.length > 0) {
    const event = logs[0];
    console.log("\n--- Event Emitted ---");
    console.log("Event name:", event?.name);
    console.log("From:", event?.args[0]);
    console.log("To:", event?.args[1]);
    console.log("Owner:", event?.args[2]);
    console.log("Actref:", event?.args[3]);
    console.log("Topic:", event?.args[4]);
  }

  // Demonstrate multiple data emissions
  console.log("\n--- Multiple Data Emissions ---");
  const ownerParam2 = ethers.keccak256(ethers.toUtf8Bytes("batch_process_001"));
  const actref2 = ethers.keccak256(ethers.toUtf8Bytes("verification_token"));
  const topic2 = "Batch Processing";

  console.log("Sending data to User2:", user2.address);
  await adminContract.connect(admin).sendDataToTarget(user2.address, ownerParam2, actref2, topic2);

  // Send data with zero values
  console.log("\n--- Sending Zero Data ---");
  const zeroData = "0x0000000000000000000000000000000000000000000000000000000000000000";
  const emptyTopic = "";
  console.log("Sending zero data to User1");
  await adminContract.connect(admin).sendDataToTarget(user1.address, zeroData, zeroData, emptyTopic);

  // Demonstrate access control (this should fail)
  console.log("\n--- Access Control Test ---");
  try {
    await adminContract.connect(user1).sendDataToTarget(user2.address, ownerParam, actref, topic);
    console.log("ERROR: Non-admin was able to send data!");
  } catch (error: any) {
    console.log("✓ Access control working: Non-admin cannot send data");
    console.log("Error:", error.message.split("(")[0]);
  }

  // Demonstrate different topic lengths
  console.log("\n--- Testing Different Topic Lengths ---");
  const shortTopic = "Hi";
  const longTopic = "This is a very long topic string that demonstrates the contract can handle various string lengths without issues and can store meaningful metadata about transactions";
  
  console.log("Short topic test...");
  await adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam, actref, shortTopic);
  
  console.log("Long topic test...");
  await adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam, actref, longTopic);

  console.log("\n=== Demo Complete ===");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
