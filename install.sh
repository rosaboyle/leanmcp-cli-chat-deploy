#!/bin/bash

# LeanMCP CLI Installation Script

set -e

echo "🚀 Installing LeanMCP CLI..."

# Build the CLI
echo "📦 Building CLI..."
go build -o leanmcp-cli .

# Make it executable
chmod +x leanmcp-cli

# Create a symlink for the alias (optional)
if command -v leanmcp-cli &> /dev/null; then
    echo "⚠️  leanmcp-cli is already installed in PATH"
else
    echo "📋 CLI built successfully as ./leanmcp-cli"
fi

echo ""
echo "✅ Installation complete!"
echo ""
echo "📖 Getting Started:"
echo "1. Authenticate with your API key:"
echo "   ./leanmcp-cli auth login --api-key airtrain_your_key_here"
echo ""
echo "2. List your projects:"
echo "   ./leanmcp-cli projects list"
echo ""
echo "3. Create an alias (optional):"
echo "   echo 'alias lcli=\"$(pwd)/leanmcp-cli\"' >> ~/.zshrc"
echo "   source ~/.zshrc"
echo ""
echo "4. Get help:"
echo "   ./leanmcp-cli --help"
echo ""
echo "🎉 Happy building with LeanMCP CLI!"
