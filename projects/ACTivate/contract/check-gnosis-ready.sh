#!/bin/bash

echo "🔍 Gnosis Chain Deployment Readiness Check"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check functions
check_file() {
    if [ -f "$1" ]; then
        echo -e "  ✅ $1 exists"
        return 0
    else
        echo -e "  ❌ $1 missing"
        return 1
    fi
}

check_env_var() {
    if grep -q "^$1=" .env 2>/dev/null && [ -n "$(grep "^$1=" .env | cut -d'=' -f2)" ]; then
        echo -e "  ✅ $1 is set"
        return 0
    else
        echo -e "  ❌ $1 not set or empty"
        return 1
    fi
}

# Initialize status
all_good=true

echo "✅ Environment files:"
check_file ".env.example" || all_good=false
check_file ".env" || all_good=false

echo ""
echo "✅ Contract files:"
check_file "contracts/AdminContract.sol" || all_good=false
check_file "scripts/deployToGnosis.ts" || all_good=false
check_file "scripts/deployToChiado.ts" || all_good=false

echo ""
echo "✅ Environment variables:"
check_env_var "PRIVATE_KEY" || all_good=false
check_env_var "GNOSIS_RPC_URL" || all_good=false
check_env_var "CHIADO_RPC_URL" || all_good=false

echo ""
echo "✅ Dependencies:"
if [ -d "node_modules" ]; then
    echo "  ✅ Dependencies installed"
else
    echo "  ❌ Dependencies not installed - run: npm install"
    all_good=false
fi

if [ -d "artifacts/contracts" ]; then
    echo "  ✅ Contracts compiled"
else
    echo "  ❌ Contracts not compiled - run: npm run compile"
    all_good=false
fi

echo ""
echo "✅ Go integration:"
check_file "examples/go/gnosis-interaction.go" || all_good=false
check_file "examples/go/go.mod" || all_good=false

echo ""
if [ "$all_good" = true ]; then
    echo -e "${GREEN}🚀 Ready to deploy!${NC}"
    echo ""
    echo "Available deployment commands:"
    echo "  npm run deploy:gnosis   # Deploy to Gnosis Chain mainnet"
    echo "  npm run deploy:chiado   # Deploy to Chiado testnet"
    echo ""
    echo "Other useful commands:"
    echo "  npm run balance:gnosis  # Check xDAI balance"
    echo "  npm run verify:gnosis   # Verify contract on Gnosisscan"
    echo "  npm run verify:chiado   # Verify contract on Chiado explorer"
    echo ""
    echo "⚠️  Remember to fund your address with xDAI first!"
    echo "   - Gnosis mainnet: Bridge at https://bridge.gnosischain.com/"
    echo "   - Chiado testnet: Get testnet xDAI at https://gnosisfaucet.com/"
else
    echo -e "${RED}❌ Not ready to deploy${NC}"
    echo ""
    echo "Please fix the issues above before deploying."
fi
