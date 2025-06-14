// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title DataContract
 * @dev Contract with public data emission functionality
 * Anyone can emit events with target addresses and custom data
 */
contract DataContract {
    // Event emitted when someone sends data to target
    event DataSentToTarget(
        address indexed from,
        address indexed to,
        bytes32 owner,
        bytes32 actref,
        string topic
    );

    /**
     * @dev Constructor - no special initialization needed
     */
    constructor() {
        // No admin setup needed
    }

    /**
     * @dev Public function to emit data to a target address
     * Anyone can call this function
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
    ) external {
        require(target != address(0), "DataContract: target cannot be zero address");

        // Emit event with the data (from = msg.sender, the actual caller)
        emit DataSentToTarget(msg.sender, target, ownerParam, actref, topic);
    }
}
