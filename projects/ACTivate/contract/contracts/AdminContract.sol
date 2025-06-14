// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title AdminContract
 * @dev Contract with admin-controlled data emission functionality
 * The contract admin (deployer) can emit events with target addresses and custom data
 */
contract AdminContract is Ownable {
    // Event emitted when admin sends data to target
    event DataSentToTarget(
        address indexed from,
        address indexed to,
        bytes32 owner,
        bytes32 actref,
        string topic
    );

    /**
     * @dev Constructor that sets the deployer as admin
     */
    constructor() Ownable(msg.sender) {
        // Admin is set by Ownable constructor
    }

    /**
     * @dev Admin function to emit data to a target address
     * Only the contract owner (admin/deployer) can call this function
     * @param target The target address
     * @param ownerParam First 64-byte (32-byte) data parameter representing owner
     * @param actref Second 64-byte (32-byte) data parameter representing action reference
     * @param topic String parameter for the topic
     */
    function sendDataToTarget(
        address target,
        bytes32 ownerParam,
        bytes32 actref,
        string calldata topic
    ) external onlyOwner {
        require(target != address(0), "AdminContract: target cannot be zero address");

        // Emit event with the data
        emit DataSentToTarget(owner(), target, ownerParam, actref, topic);
    }

    /**
     * @dev Get the admin address
     * @return The address of the contract admin/owner
     */
    function getAdmin() external view returns (address) {
        return owner();
    }
}
