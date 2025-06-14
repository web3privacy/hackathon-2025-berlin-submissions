import { expect } from "chai";
import { ethers } from "hardhat";
import { AdminContract } from "../typechain-types";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";

describe("AdminContract", function () {
  let adminContract: AdminContract;
  let admin: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let user2: HardhatEthersSigner;

  beforeEach(async function () {
    // Get signers
    [admin, user1, user2] = await ethers.getSigners();

    // Deploy the contract
    const AdminContractFactory = await ethers.getContractFactory("AdminContract");
    adminContract = await AdminContractFactory.deploy();
    await adminContract.waitForDeployment();
  });

  describe("Deployment", function () {
    it("Should set the right admin", async function () {
      expect(await adminContract.getAdmin()).to.equal(admin.address);
      expect(await adminContract.owner()).to.equal(admin.address);
    });
  });

  describe("Admin Functions", function () {
    describe("sendDataToTarget", function () {
      it("Should allow admin to send data to target", async function () {
        const ownerParam = ethers.keccak256(ethers.toUtf8Bytes("owner"));
        const actref = ethers.keccak256(ethers.toUtf8Bytes("actref"));
        const topic = "Test Topic";

        await expect(
          adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam, actref, topic)
        )
          .to.emit(adminContract, "DataSentToTarget")
          .withArgs(admin.address, user1.address, ownerParam, actref, topic);
      });

      it("Should revert if non-admin tries to send data", async function () {
        const ownerParam = ethers.keccak256(ethers.toUtf8Bytes("owner"));
        const actref = ethers.keccak256(ethers.toUtf8Bytes("actref"));
        const topic = "Test Topic";

        await expect(
          adminContract.connect(user1).sendDataToTarget(user2.address, ownerParam, actref, topic)
        ).to.be.revertedWithCustomError(adminContract, "OwnableUnauthorizedAccount");
      });

      it("Should revert if target is zero address", async function () {
        const ownerParam = ethers.keccak256(ethers.toUtf8Bytes("owner"));
        const actref = ethers.keccak256(ethers.toUtf8Bytes("actref"));
        const topic = "Test Topic";

        await expect(
          adminContract.connect(admin).sendDataToTarget(ethers.ZeroAddress, ownerParam, actref, topic)
        ).to.be.revertedWith("AdminContract: target cannot be zero address");
      });

      it("Should handle multiple data emissions correctly", async function () {
        const ownerParam1 = ethers.keccak256(ethers.toUtf8Bytes("owner1"));
        const actref1 = ethers.keccak256(ethers.toUtf8Bytes("actref1"));
        const topic1 = "Topic 1";
        
        const ownerParam2 = ethers.keccak256(ethers.toUtf8Bytes("owner2"));
        const actref2 = ethers.keccak256(ethers.toUtf8Bytes("actref2"));
        const topic2 = "Topic 2";

        // First emission
        await expect(
          adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam1, actref1, topic1)
        )
          .to.emit(adminContract, "DataSentToTarget")
          .withArgs(admin.address, user1.address, ownerParam1, actref1, topic1);

        // Second emission
        await expect(
          adminContract.connect(admin).sendDataToTarget(user2.address, ownerParam2, actref2, topic2)
        )
          .to.emit(adminContract, "DataSentToTarget")
          .withArgs(admin.address, user2.address, ownerParam2, actref2, topic2);
      });
    });
  });

  describe("Data Handling", function () {
    it("Should correctly handle different data formats", async function () {
      // Test with hex data
      const ownerParam = "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
      const actref = "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321";
      const topic = "Hex Data Topic";

      await expect(
        adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam, actref, topic)
      )
        .to.emit(adminContract, "DataSentToTarget")
        .withArgs(admin.address, user1.address, ownerParam, actref, topic);
    });

    it("Should handle zero data values", async function () {
      const zeroData = "0x0000000000000000000000000000000000000000000000000000000000000000";
      const topic = "Zero Data Topic";

      await expect(
        adminContract.connect(admin).sendDataToTarget(user1.address, zeroData, zeroData, topic)
      )
        .to.emit(adminContract, "DataSentToTarget")
        .withArgs(admin.address, user1.address, zeroData, zeroData, topic);
    });

    it("Should handle same target multiple times", async function () {
      const ownerParam1 = ethers.keccak256(ethers.toUtf8Bytes("first_owner"));
      const actref1 = ethers.keccak256(ethers.toUtf8Bytes("first_actref"));
      const topic1 = "First Topic";
      
      const ownerParam2 = ethers.keccak256(ethers.toUtf8Bytes("second_owner"));
      const actref2 = ethers.keccak256(ethers.toUtf8Bytes("second_actref"));
      const topic2 = "Second Topic";

      // Send to same target twice
      await expect(
        adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam1, actref1, topic1)
      )
        .to.emit(adminContract, "DataSentToTarget")
        .withArgs(admin.address, user1.address, ownerParam1, actref1, topic1);

      await expect(
        adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam2, actref2, topic2)
      )
        .to.emit(adminContract, "DataSentToTarget")
        .withArgs(admin.address, user1.address, ownerParam2, actref2, topic2);
    });

    it("Should handle empty topic string", async function () {
      const ownerParam = ethers.keccak256(ethers.toUtf8Bytes("owner"));
      const actref = ethers.keccak256(ethers.toUtf8Bytes("actref"));
      const topic = "";

      await expect(
        adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam, actref, topic)
      )
        .to.emit(adminContract, "DataSentToTarget")
        .withArgs(admin.address, user1.address, ownerParam, actref, topic);
    });

    it("Should handle long topic string", async function () {
      const ownerParam = ethers.keccak256(ethers.toUtf8Bytes("owner"));
      const actref = ethers.keccak256(ethers.toUtf8Bytes("actref"));
      const topic = "This is a very long topic string that contains multiple words and demonstrates that the contract can handle longer string parameters without any issues";

      await expect(
        adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam, actref, topic)
      )
        .to.emit(adminContract, "DataSentToTarget")
        .withArgs(admin.address, user1.address, ownerParam, actref, topic);
    });
  });

  describe("Access Control", function () {
    it("Should return correct admin address", async function () {
      expect(await adminContract.getAdmin()).to.equal(admin.address);
    });

    it("Should maintain admin privileges", async function () {
      const ownerParam = ethers.keccak256(ethers.toUtf8Bytes("test"));
      const actref = ethers.keccak256(ethers.toUtf8Bytes("test"));
      const topic = "Test Topic";

      // Admin should be able to call function
      await expect(
        adminContract.connect(admin).sendDataToTarget(user1.address, ownerParam, actref, topic)
      ).to.not.be.reverted;

      // Non-admin should not be able to call function
      await expect(
        adminContract.connect(user1).sendDataToTarget(user2.address, ownerParam, actref, topic)
      ).to.be.revertedWithCustomError(adminContract, "OwnableUnauthorizedAccount");
    });
  });
});
