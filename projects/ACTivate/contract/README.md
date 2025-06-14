# AdminContract

A secure smart contract with admin-controlled data emission functionality, built with Hardhat and OpenZeppelin.

## Features

- **Admin Control**: Only the contract deployer can emit events with custom data
- **Secure Access**: Uses OpenZeppelin's `Ownable` for access control
- **Structured Data**: Emit events with owner, action reference, and topic parameters
- **Event Logging**: Emits `DataSentToTarget` events for comprehensive tracking
- **No Token Functionality**: Pure event emission contract without any token transfers

## Contract Details

- **Contract Name**: AdminContract
- **Admin**: Contract deployer address
- **Purpose**: Emit events with structured data to target addresses

## Event Structure

The contract emits `DataSentToTarget` events with the following structure:
```solidity
event DataSentToTarget(
    address indexed from,    // Admin address (event sender)
    address indexed to,      // Target address
    bytes32 owner,          // Owner identifier (32 bytes)
    bytes32 actref,         // Action reference (32 bytes)
    string topic            // Topic description
);
```

## Key Functions

### `sendDataToTarget(address target, bytes32 ownerParam, bytes32 actref, string calldata topic)`
- **Access**: Admin only
- **Purpose**: Emit structured data to a target address
- **Parameters**:
  - `target`: Target address
  - `ownerParam`: Owner identifier (32 bytes)
  - `actref`: Action reference identifier (32 bytes)
  - `topic`: Topic description (string)
- **Events**: Emits `DataSentToTarget` event

### `getAdmin()`
- **Access**: Public view
- **Purpose**: Returns the admin address

## Development Setup

### Prerequisites
- Node.js (v16+ recommended)
- npm or yarn

### Installation

1. Clone the repository
2. Install dependencies:
```bash
npm install
```

### Available Scripts

```bash
# Compile contracts
npm run compile

# Run tests
npm run test

# Deploy to local network
npm run deploy

# Deploy to Sepolia testnet
npm run deploy:sepolia

# Run interaction demo
npm run interact

# Verify contract on Sepolia
npm run verify:sepolia CONTRACT_ADDRESS

# Check Sepolia deployment readiness
npm run check:sepolia

# Run complete test suite
npm run test:complete

# Start local Hardhat node
npm run node

# Clean artifacts
npm run clean
```

## Sepolia Testnet Deployment

Deploy to Sepolia testnet with full Go integration support:

### Quick Setup
1. Copy environment template: `cp .env.example .env`
2. Add your Sepolia RPC URL and private key to `.env`
3. Get testnet ETH from [Sepolia Faucet](https://sepoliafaucet.com/)
4. Deploy: `npm run deploy:sepolia`

The deployment script automatically generates:
- ✅ **JSON configuration** for general use
- ✅ **Go constants file** ready to copy to your Go project  
- ✅ **Environment variables** for easy loading
- ✅ **Complete Go configuration struct** with embedded ABI

See [SEPOLIA_DEPLOYMENT.md](./SEPOLIA_DEPLOYMENT.md) for detailed instructions.

## Helper Tools

### Deployment Readiness Check
```bash
npm run check:sepolia
```
Validates that your environment is properly configured for Sepolia deployment.

### Complete Test Suite
```bash
npm run test:complete
```
Runs a comprehensive test of the entire project including:
- Contract compilation and testing
- Local deployment validation
- TypeScript compilation checks
- Go integration file validation
- Documentation verification
- Security checks

## Testing

The contract includes comprehensive tests covering:
- Deployment verification
- Admin functionality
- Access control
- Data emission
- Error conditions

Run tests with:
```bash
npm run test
```

**Test Results: 12 passing tests**
- ✅ Should set the right admin
- ✅ Should allow admin to send data to target
- ✅ Should revert if non-admin tries to send data
- ✅ Should revert if target is zero address
- ✅ Should handle multiple data emissions correctly
- ✅ Should correctly handle different data formats
- ✅ Should handle zero data values
- ✅ Should handle same target multiple times
- ✅ Should return correct admin address
- ✅ Should maintain admin privileges
- ✅ Should emit correct event data
- ✅ Should handle string topics correctly

## Deployment

### Local Development

1. Start a local Hardhat node:
```bash
npm run node
```

2. Deploy the contract:
```bash
npm run deploy
```

### Sepolia Testnet Deployment

1. **Setup Environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your Sepolia RPC URL and private key
   ```

2. **Get Testnet ETH**:
   Visit [Sepolia Faucet](https://sepoliafaucet.com/) to get free testnet ETH

3. **Deploy Contract**:
   ```bash
   npm run deploy:sepolia
   ```

4. **Verify Contract** (optional):
   ```bash
   npm run verify:sepolia CONTRACT_ADDRESS
   ```

The deployment script generates multiple output formats:
- `deployments/sepolia-deployment.json` - Complete deployment info
- `deployments/constants.go` - Go constants ready to use
- `deployments/sepolia.env` - Environment variables
- `deployments/sepolia-config.go` - Go configuration struct

### Production Deployment

Update `hardhat.config.ts` with your target network configuration and run:
```bash
npm run deploy --network <network-name>
```

## Security Features

- **Access Control**: Uses OpenZeppelin's battle-tested `Ownable` pattern
- **Input Validation**: Validates target addresses (no zero address)
- **Event Logging**: All admin actions are logged via events
- **Gas Efficient**: Minimal contract functionality for reduced gas costs

## Usage Example

```solidity
// Deploy the contract
AdminContract contract = new AdminContract();

// Send data to target with structured parameters
bytes32 ownerParam = keccak256("OWNER_001");
bytes32 actref = keccak256("ACTION_REF_12345");
string memory topic = "Data Transfer Event";
contract.sendDataToTarget(targetAddress, ownerParam, actref, topic);
```

## Events

### DataSentToTarget
Emitted when admin sends data to target:
```solidity
event DataSentToTarget(
    address indexed from,    // Admin address
    address indexed to,      // Target address
    bytes32 owner,          // Owner identifier
    bytes32 actref,         // Action reference
    string topic            // Topic description
);
```

## Examples

### TypeScript Interaction
See `scripts/interactContract.ts` for a complete TypeScript example demonstrating:
- Contract deployment and connection
- Admin function calls with structured data
- Event parsing and logging
- Access control testing
- Batch data operations

### Go Integration
See `examples/go/` for a complete Go example demonstrating:
- Ethereum client connection
- Contract interaction using go-ethereum
- Transaction handling and receipt parsing
- Event filtering and historical queries
- Real-world blockchain integration patterns

## Interaction Demo Results

The interaction demo successfully demonstrates:
- Contract deployment with admin privileges
- Sending structured data (owner, action reference, topic) to target addresses
- Event emission with comprehensive data structure
- Access control preventing unauthorized data emission
- Handling of various data formats and string topics
- Multiple data emissions to different targets
- Event filtering and historical data retrieval

## License

MIT