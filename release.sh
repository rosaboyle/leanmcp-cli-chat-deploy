#!/bin/bash

# GitHub Release Script for LeanMCP
# Usage: ./release.sh <version>
# Example: ./release.sh 1.0.0

set -e

VERSION=${1:-"1.0.0"}
REPO_OWNER="rosaboyle"  # Your GitHub username
REPO_NAME="leanmcp-cli-chat-deploy"  # Your existing repo name

echo " Creating release for leanmcp v${VERSION}"

# Check if version is provided
if [ -z "$1" ]; then
    echo "⚠️  No version provided, using default: ${VERSION}"
    echo "Usage: ./release.sh <version>"
    echo "Example: ./release.sh 1.0.0"
fi

# Clean and build
echo "🧹 Cleaning previous builds..."
make clean

echo "🔨 Building release binaries..."
make release VERSION=${VERSION}

# Check if binaries were created
if [ ! -f "leanmcp-${VERSION}-darwin-amd64.tar.gz" ] || [ ! -f "leanmcp-${VERSION}-darwin-arm64.tar.gz" ]; then
    echo "❌ Release binaries not found. Build may have failed."
    exit 1
fi

# Calculate checksums
echo "🔢 Calculating checksums..."
SHA256_AMD64=$(shasum -a 256 leanmcp-${VERSION}-darwin-amd64.tar.gz | cut -d' ' -f1)
SHA256_ARM64=$(shasum -a 256 leanmcp-${VERSION}-darwin-arm64.tar.gz | cut -d' ' -f1)

echo "✅ Release artifacts created:"
echo "📦 leanmcp-${VERSION}-darwin-amd64.tar.gz (SHA256: ${SHA256_AMD64})"
echo "📦 leanmcp-${VERSION}-darwin-arm64.tar.gz (SHA256: ${SHA256_ARM64})"

echo ""
echo "🍺 Homebrew Formula Template:"
echo "============================================"
cat << EOF
class Leanmcp < Formula
  desc "LeanMCP CLI - Manage projects and chats"
  homepage "https://github.com/${REPO_OWNER}/${REPO_NAME}"
  version "${VERSION}"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/v${VERSION}/leanmcp-${VERSION}-darwin-arm64.tar.gz"
    sha256 "${SHA256_ARM64}"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/v${VERSION}/leanmcp-${VERSION}-darwin-amd64.tar.gz"
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
echo "============================================"

echo ""
echo "📋 Next Steps:"
echo "1. Create a GitHub release v${VERSION}"
echo "2. Upload these files to the release:"
echo "   - leanmcp-${VERSION}-darwin-amd64.tar.gz"
echo "   - leanmcp-${VERSION}-darwin-arm64.tar.gz"
echo "3. Copy the Homebrew formula above to your tap repository"
echo "4. Update REPO_OWNER and REPO_NAME in this script if needed"

echo ""
echo "🔗 GitHub Release URL:"
echo "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/new?tag=v${VERSION}&title=v${VERSION}"
