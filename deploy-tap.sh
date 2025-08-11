#!/bin/bash

# Deploy Homebrew Tap Script
# Run this AFTER creating the homebrew-mcli repository on GitHub

set -e

echo "🍺 Deploying Homebrew Tap..."

# Clone the tap repository
if [ -d "homebrew-mcli" ]; then
    echo "📁 Directory homebrew-mcli exists, removing..."
    rm -rf homebrew-mcli
fi

echo "📥 Cloning homebrew-mcli repository..."
git clone https://github.com/rosaboyle/homebrew-mcli.git

# Copy files
echo "📝 Adding formula and README..."
cp mcli.rb homebrew-mcli/
cp TAP_README.md homebrew-mcli/README.md

# Push to repository
cd homebrew-mcli
git add .
git config user.name "rosaboyle" 2>/dev/null || echo "Git user already configured"
git config user.email "helloworldcmu@gmail.com" 2>/dev/null || echo "Git email already configured"
git commit -m "Add mcli formula v$(grep 'version' ../mcli.rb | sed 's/.*"\(.*\)".*/\1/')"
git push

echo "✅ Homebrew tap deployed successfully!"
echo ""
echo "🎉 Users can now install MCLI with:"
echo "brew tap rosaboyle/mcli"
echo "brew install mcli"

cd ..
rm -rf homebrew-mcli
