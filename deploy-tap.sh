#!/bin/bash

# Deploy Homebrew Tap Script
# Run this AFTER creating the homebrew-leanmcp repository on GitHub

set -e

echo "🍺 Deploying Homebrew Tap..."

# Clone the tap repository
if [ -d "homebrew-leanmcp" ]; then
    echo "📁 Directory homebrew-leanmcp exists, removing..."
    rm -rf homebrew-leanmcp
fi

echo "📥 Cloning homebrew-leanmcp repository..."
git clone https://github.com/rosaboyle/homebrew-leanmcp.git

# Copy files
echo "📝 Adding formula and README..."
cp leanmcp.rb homebrew-leanmcp/
cp TAP_README.md homebrew-leanmcp/README.md

# Push to repository
cd homebrew-leanmcp
git add .
git config user.name "rosaboyle" 2>/dev/null || echo "Git user already configured"
git config user.email "helloworldcmu@gmail.com" 2>/dev/null || echo "Git email already configured"
git commit -m "Add leanmcp formula v$(grep 'version' ../leanmcp.rb | sed 's/.*"\(.*\)".*/\1/')"
git push

echo "✅ Homebrew tap deployed successfully!"
echo ""
echo "🎉 Users can now install LeanMCP with:"
echo "brew tap rosaboyle/leanmcp"
echo "brew install leanmcp"

cd ..
rm -rf homebrew-leanmcp
