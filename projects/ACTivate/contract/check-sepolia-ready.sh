#!/bin/bash
echo "🔍 Sepolia Deployment Readiness Check"
echo "======================================"
echo ""

# Quick checks
echo "✅ Environment files:"
[ -f .env.example ] && echo "  ✅ .env.example exists" || echo "  ❌ .env.example missing"
[ -f .env ] && echo "  ✅ .env exists" || echo "  ⚠️  .env missing (copy from .env.example)"

echo ""
echo "✅ Contract files:"
[ -f contracts/AdminContract.sol ] && echo "  ✅ AdminContract.sol exists" || echo "  ❌ AdminContract.sol missing"
[ -f scripts/deployToSepolia.ts ] && echo "  ✅ deployToSepolia.ts exists" || echo "  ❌ deployToSepolia.ts missing"

echo ""
echo "✅ Dependencies:"
[ -d node_modules ] && echo "  ✅ Dependencies installed" || echo "  ❌ Run: npm install"
[ -d artifacts ] && echo "  ✅ Contracts compiled" || echo "  ❌ Run: npm run compile"

echo ""
echo "✅ Go integration:"
[ -f examples/go/sepolia-interaction.go ] && echo "  ✅ Go example exists" || echo "  ❌ Go example missing"

echo ""
echo "🚀 Ready to deploy with: npm run deploy:sepolia"
