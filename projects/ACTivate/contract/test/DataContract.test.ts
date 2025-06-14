import { expect } from "chai";
import { ethers } from "hardhat";
import { DataContract } from "../typechain-types";
import { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers";

describe("DataContract", function () {
  let dataContract: DataContract;
  let owner: SignerWithAddress;
  let user1: SignerWithAddress;
  let user2: SignerWithAddress;
  let targetAddress: SignerWithAddress;

  beforeEach(async function () {
    // Get test accounts
    [owner, user1, user2, targetAddress] = await ethers.getSigners();

    // Deploy the contract
    const DataContract = await ethers.getContractFactory("DataContract");
    dataContract = await DataContract.deploy();
    await dataContract.waitForDeployment();
  });

  describe("Deployment", function () {
    it("Should deploy successfully", async function () {
      expect(await dataContract.getAddress()).to.be.properAddress;
    });
  });

  describe("Public Functions", function () {
    describe("sendDataToTarget", function () {
      it("Should allow anyone to send data to target", async function () {
        const ownerParam = ethers.encodeBytes32String("OWNER_001");
        const actref = ethers.encodeBytes32String("ACTION_REF_123");
        const topic = "Test Topic";

        // Test with contract deployer
        await expect(
          dataContract.sendDataToTarget(targetAddress.address, ownerParam, actref, topic)
        )
          .to.emit(dataContract, "DataSentToTarget")
          .withArgs(owner.address, targetAddress.address, ownerParam, actref, topic);

        // Test with user1 (should work - public access)
        await expect(
          dataContract.connect(user1).sendDataToTarget(targetAddress.address, ownerParam, actref, topic)
        )
          .to.emit(dataContract, "DataSentToTarget")
          .withArgs(user1.address, targetAddress.address, ownerParam, actref, topic);

        // Test with user2 (should work - public access)
        await expect(
          dataContract.connect(user2).sendDataToTarget(targetAddress.address, ownerParam, actref, topic)
        )
          .to.emit(dataContract, "DataSentToTarget")
          .withArgs(user2.address, targetAddress.address, ownerParam, actref, topic);
      });

      it("Should revert if target is zero address", async function () {
        const ownerParam = ethers.encodeBytes32String("OWNER_001");
        const actref = ethers.encodeBytes32String("ACTION_REF_123");
        const topic = "Test Topic";

        await expect(
          dataContract.sendDataToTarget(ethers.ZeroAddress, ownerParam, actref, topic)
        ).to.be.revertedWith("DataContract: target cannot be zero address");
      });

      it("Should handle multiple data emissions correctly", async function () {
        const ownerParam1 = ethers.encodeBytes32String("OWNER_001");
        const ownerParam2 = ethers.encodeBytes32String("OWNER_002");
        const actref1 = ethers.encodeBytes32String("ACTION_REF_123");
        const actref2 = ethers.encodeBytes32String("ACTION_REF_456");
        const topic1 = "First Topic";
        const topic2 = "Second Topic";

        // First emission from owner
        await expect(
          dataContract.sendDataToTarget(targetAddress.address, ownerParam1, actref1, topic1)
        )
          .to.emit(dataContract, "DataSentToTarget")
          .withArgs(owner.address, targetAddress.address, ownerParam1, actref1, topic1);

        // Second emission from user1
        await expect(
          dataContract.connect(user1).sendDataToTarget(user2.address, ownerParam2, actref2, topic2)
        )
          .to.emit(dataContract, "DataSentToTarget")
          .withArgs(user1.address, user2.address, ownerParam2, actref2, topic2);
      });

      it("Should emit correct event data", async function () {
        const ownerParam = ethers.encodeBytes32String("TEST_OWNER");
        const actref = ethers.encodeBytes32String("TEST_ACTION");
        const topic = "Test Event Data";

        const tx = await dataContract.connect(user1).sendDataToTarget(
          targetAddress.address,
          ownerParam,
          actref,
          topic
        );

        const receipt = await tx.wait();
        expect(receipt).to.not.be.null;

        // Verify event was emitted
        await expect(tx)
          .to.emit(dataContract, "DataSentToTarget")
          .withArgs(user1.address, targetAddress.address, ownerParam, actref, topic);
      });
    });
  });

  describe("Multi-User Access", function () {
    it("Should allow multiple users to call the function simultaneously", async function () {
      const ownerParam = ethers.encodeBytes32String("MULTI_USER");
      const actref = ethers.encodeBytes32String("PARALLEL_TEST");
      const topic = "Multi-user test";

      // All users should be able to call the function
      const promises = [
        dataContract.connect(owner).sendDataToTarget(targetAddress.address, ownerParam, actref, topic + " - Owner"),
        dataContract.connect(user1).sendDataToTarget(targetAddress.address, ownerParam, actref, topic + " - User1"),
        dataContract.connect(user2).sendDataToTarget(targetAddress.address, ownerParam, actref, topic + " - User2"),
      ];

      // All transactions should succeed
      const results = await Promise.all(promises);
      expect(results).to.have.length(3);

      // All should have transaction receipts
      for (const result of results) {
        const receipt = await result.wait();
        expect(receipt).to.not.be.null;
        expect(receipt?.blockNumber).to.be.greaterThan(0);
      }
    });
  });
});