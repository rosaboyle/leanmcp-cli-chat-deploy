#!/bin/bash

# Homebrew Tap Setup Script for LeanMCP
# This creates the homebrew-leanmcp repository and formula

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🍺 Setting up Homebrew Tap for LeanMCP${NC}"
echo "=================================="

# Get the latest release info from GitHub API
echo -e "${BLUE}📡 Fetching latest release info...${NC}"
LATEST_RELEASE=$(curl -s https://api.github.com/repos/rosaboyle/leanmcp-cli-chat-deploy/releases/latest)

if [ "$(echo $LATEST_RELEASE | grep '\"message\": \"Not Found\"')" ]; then
    echo -e "${RED}❌ No releases found. Please run ./auto-release.sh first${NC}"
    exit 1
fi

VERSION=$(echo $LATEST_RELEASE | grep '"tag_name"' | sed -E 's/.*"tag_name": "v([^"]+)".*/\1/')
AMD64_URL=$(echo $LATEST_RELEASE | grep '"browser_download_url"' | grep 'darwin-amd64' | sed -E 's/.*"browser_download_url": "([^"]+)".*/\1/')
ARM64_URL=$(echo $LATEST_RELEASE | grep '"browser_download_url"' | grep 'darwin-arm64' | sed -E 's/.*"browser_download_url": "([^"]+)".*/\1/')

if [ -z "$VERSION" ] || [ -z "$AMD64_URL" ] || [ -z "$ARM64_URL" ]; then
    echo -e "${RED}❌ Could not extract release info. Please check the release exists:${NC}"
    echo "https://github.com/rosaboyle/leanmcp-cli-chat-deploy/releases"
    exit 1
fi

echo -e "${GREEN}✅ Found release v${VERSION}${NC}"

# Download and calculate checksums
echo -e "${BLUE}🔢 Calculating checksums...${NC}"
TMP_DIR=$(mktemp -d)
curl -sL "$AMD64_URL" -o "$TMP_DIR/amd64.tar.gz"
curl -sL "$ARM64_URL" -o "$TMP_DIR/arm64.tar.gz"

SHA256_AMD64=$(shasum -a 256 "$TMP_DIR/amd64.tar.gz" | cut -d' ' -f1)
SHA256_ARM64=$(shasum -a 256 "$TMP_DIR/arm64.tar.gz" | cut -d' ' -f1)

rm -rf "$TMP_DIR"

echo -e "${GREEN}✅ Checksums calculated${NC}"

# Create formula
echo -e "${BLUE}📝 Generating Homebrew formula...${NC}"
cat > leanmcp.rb << EOF
class Leanmcp < Formula
  desc "LeanMCP CLI - Manage projects and chats"
  homepage "https://github.com/rosaboyle/leanmcp-cli-chat-deploy"
  version "${VERSION}"

  if OS.mac? && Hardware::CPU.arm?
    url "${ARM64_URL}"
    sha256 "${SHA256_ARM64}"
  elsif OS.mac? && Hardware::CPU.intel?
    url "${AMD64_URL}"
    sha256 "${SHA256_AMD64}"
  end

  def install
    bin.install "leanmcp-darwin-arm64" => "leanmcp" if OS.mac? && Hardware::CPU.arm?
    bin.install "leanmcp-darwin-amd64" => "leanmcp" if OS.mac? && Hardware::CPU.intel?
  end

  test do
    system "#{bin}/leanmcp", "version"
  end
end
EOF

echo -e "${GREEN}✅ Formula created: leanmcp.rb${NC}"

# Create README
cat > TAP_README.md << 'EOF'
# Homebrew Tap for LeanMCP

This is the official Homebrew tap for LeanMCP CLI.

## Installation

```bash
brew tap rosaboyle/leanmcp
brew install leanmcp
```

## Usage

```bash
leanmcp --help
leanmcp version
leanmcp projects list
leanmcp chats list
```

## Repository

Source code: https://github.com/rosaboyle/leanmcp-cli-chat-deploy
EOF

echo ""
echo -e "${YELLOW}🚀 Next Steps:${NC}"
echo "1. Create a new GitHub repository named: ${BLUE}homebrew-leanmcp${NC}"
echo "2. Clone it and add these files:"
echo "   - leanmcp.rb (generated above)"
echo "   - README.md (generated as TAP_README.md)"
echo ""
echo -e "${BLUE}📋 Quick setup commands:${NC}"
echo "# After creating the GitHub repo homebrew-leanmcp:"
echo "git clone https://github.com/rosaboyle/homebrew-leanmcp.git"
echo "cd homebrew-leanmcp"
echo "cp $(pwd)/leanmcp.rb ."
echo "cp $(pwd)/TAP_README.md README.md"
echo "git add ."
echo 'git commit -m "Add leanmcp formula"'
echo "git push"
echo ""
echo -e "${GREEN}🎉 Then users can install with:${NC}"
echo "brew tap rosaboyle/leanmcp && brew install leanmcp"
echo ""
echo -e "${BLUE}📁 Files created in current directory:${NC}"
echo "- leanmcp.rb (Homebrew formula)"
echo "- TAP_README.md (README for tap repo)"
