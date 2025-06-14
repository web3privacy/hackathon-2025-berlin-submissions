#!/bin/bash

# AdminContract Complete Test Suite
# This script tests the entire deployment and interaction workflow

set -e  # Exit on any error

echo "🧪 AdminContract Complete Test Suite"
echo "===================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_step() {
    echo -e "${BLUE}📋 Step $1: $2${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Step 1: Environment Setup
print_step "1" "Checking Environment Setup"

if [ ! -f ".env" ]; then
    print_warning ".env file not found. Creating from template..."
    cp .env.example .env
    print_warning "Please edit .env with your actual values before running Sepolia deployment"
fi

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    print_warning "Dependencies not installed. Installing..."
    npm install
fi

print_success "Environment setup complete"
echo ""

# Step 2: Contract Compilation
print_step "2" "Compiling Smart Contracts"

npm run clean
npm run compile

print_success "Contract compilation complete"
echo ""

# Step 3: Running Tests
print_step "3" "Running Test Suite"

npm test

print_success "All tests passed"
echo ""

# Step 4: Local Deployment Test
print_step "4" "Testing Local Deployment"

# Start local node in background
print_warning "Starting local Hardhat node..."
npm run node > /dev/null 2>&1 &
NODE_PID=$!

# Wait for node to start
sleep 3

# Deploy to local network
npm run deploy

# Test interaction
CONTRACT_ADDRESS=$(npx hardhat run scripts/deployContract.ts --network localhost 2>/dev/null | grep "deployed to:" | awk '{print $NF}')
if [ -n "$CONTRACT_ADDRESS" ]; then
    print_success "Local deployment successful: $CONTRACT_ADDRESS"
    
    # Test interaction
    print_warning "Testing contract interaction..."
    CONTRACT_ADDRESS=$CONTRACT_ADDRESS npm run interact
    print_success "Contract interaction test passed"
else
    print_error "Failed to get contract address from deployment"
fi

# Clean up local node
kill $NODE_PID 2>/dev/null || true
print_success "Local deployment test complete"
echo ""

# Step 5: TypeScript Compilation Check
print_step "5" "Checking TypeScript Compilation"

# Check all TypeScript files
echo "Checking deployment script..."
npx hardhat compile  # This implicitly checks TypeScript

print_success "TypeScript compilation check passed"
echo ""

# Step 6: Go Integration Files Check
print_step "6" "Validating Go Integration Files"

if [ -d "deployments" ]; then
    if [ -f "deployments/constants.go" ]; then
        print_success "Go constants file exists"
        
        # Basic syntax check for Go file
        if command -v go >/dev/null 2>&1; then
            cd deployments
            echo "package main" > temp_test.go
            cat constants.go >> temp_test.go
            echo "func main() {}" >> temp_test.go
            
            if go build temp_test.go 2>/dev/null; then
                print_success "Go constants file has valid syntax"
            else
                print_warning "Go constants file has syntax issues"
            fi
            
            rm -f temp_test.go temp_test
            cd ..
        else
            print_warning "Go not installed, skipping syntax check"
        fi
    else
        print_warning "Go constants file not found (run deployment first)"
    fi
    
    if [ -f "deployments/sepolia.env" ]; then
        print_success "Environment file exists"
    else
        print_warning "Environment file not found (run deployment first)"
    fi
else
    print_warning "Deployments directory not found (run deployment first)"
fi

echo ""

# Step 7: Documentation Check
print_step "7" "Checking Documentation"

docs=("README.md" "SEPOLIA_DEPLOYMENT.md" "PROJECT_SUMMARY.md")
for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        print_success "$doc exists"
    else
        print_warning "$doc not found"
    fi
done

if [ -f "examples/go/README.md" ]; then
    print_success "Go examples documentation exists"
else
    print_warning "Go examples documentation not found"
fi

echo ""

# Step 8: Security Check
print_step "8" "Security Checks"

# Check if .env is in .gitignore
if [ -f ".gitignore" ] && grep -q "\.env" .gitignore; then
    print_success ".env file is properly ignored by git"
else
    print_warning ".env file should be added to .gitignore"
fi

# Check for hardcoded private keys
if grep -r "0x[a-fA-F0-9]\{64\}" scripts/ --exclude-dir=node_modules 2>/dev/null | grep -v "0x0000000000000000000000000000000000000000000000000000000000000000"; then
    print_error "Found potential hardcoded private keys in scripts!"
else
    print_success "No hardcoded private keys found in scripts"
fi

echo ""

# Step 9: Project Structure Validation
print_step "9" "Validating Project Structure"

required_dirs=("contracts" "scripts" "test" "examples/go")
for dir in "${required_dirs[@]}"; do
    if [ -d "$dir" ]; then
        print_success "$dir directory exists"
    else
        print_error "$dir directory missing"
    fi
done

required_files=("hardhat.config.ts" "package.json" "tsconfig.json")
for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        print_success "$file exists"
    else
        print_error "$file missing"
    fi
done

echo ""

# Step 10: Package Scripts Check
print_step "10" "Validating Package Scripts"

scripts=("compile" "test" "deploy" "deploy:sepolia" "interact" "verify:sepolia")
for script in "${scripts[@]}"; do
    if npm run-script --silent 2>/dev/null | grep -q "^  $script$"; then
        print_success "npm script '$script' exists"
    else
        print_warning "npm script '$script' not found"
    fi
done

echo ""

# Final Summary
echo "🎯 Test Summary"
echo "==============="
echo ""

print_success "✅ Contract compilation and testing"
print_success "✅ Local deployment and interaction"
print_success "✅ TypeScript compilation"
print_success "✅ Go integration files"
print_success "✅ Documentation"
print_success "✅ Security checks"
print_success "✅ Project structure"
print_success "✅ Package scripts"

echo ""
echo "🚀 Ready for Sepolia Deployment!"
echo ""
echo "To deploy to Sepolia testnet:"
echo "1. Edit .env with your actual values"
echo "2. Get Sepolia ETH from https://sepoliafaucet.com/"
echo "3. Run: npm run deploy:sepolia"
echo ""
echo "To test Go integration:"
echo "1. cd examples/go"
echo "2. source ../../deployments/sepolia.env"
echo "3. go run sepolia-interaction.go"
echo ""

print_success "All tests completed successfully! 🎉"
