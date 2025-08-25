#!/bin/bash

# Auto Release Script for LeanMCP
# Usage: ./auto-release.sh <version>
# Example: ./auto-release.sh 1.0.1

set -e

VERSION=${1}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE} LeanMCP Auto Release Script${NC}"
echo "=============================="

# Check if version is provided
if [ -z "$VERSION" ]; then
    echo -e "${RED}❌ Error: Version is required${NC}"
    echo "Usage: ./auto-release.sh <version>"
    echo "Example: ./auto-release.sh 1.0.1"
    exit 1
fi

# Validate version format (basic semver check)
if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Error: Invalid version format${NC}"
    echo "Please use semantic versioning: X.Y.Z (e.g., 1.0.1)"
    exit 1
fi

echo -e "${YELLOW}📋 Preparing release v${VERSION}${NC}"

# Check if we're on main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}❌ Error: Must be on main branch${NC}"
    echo "Current branch: $CURRENT_BRANCH"
    echo "Switch to main: git checkout main"
    exit 1
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${YELLOW}⚠️  Working directory has uncommitted changes${NC}"
    echo "Committing changes..."
    git add -A
    git commit -m "Prepare release v${VERSION}"
fi

# Check if tag already exists
if git tag -l | grep -q "v${VERSION}"; then
    echo -e "${RED}❌ Error: Tag v${VERSION} already exists${NC}"
    echo "Existing tags:"
    git tag -l | tail -5
    exit 1
fi

# Update version in root.go
echo -e "${BLUE}📝 Updating version in code...${NC}"
sed -i.bak "s/Version = \".*\"/Version = \"${VERSION}\"/" cmd/root.go
rm cmd/root.go.bak

# Commit version update
git add cmd/root.go
git commit -m "Bump version to ${VERSION}" || echo "No changes to commit"

# Pull latest changes
echo -e "${BLUE}🔄 Pulling latest changes...${NC}"
git pull origin main

# Create and push tag
echo -e "${BLUE}🏷️  Creating tag v${VERSION}...${NC}"
git tag -a "v${VERSION}" -m "Release v${VERSION}"

echo -e "${BLUE}⬆️  Pushing changes and tag...${NC}"
git push origin main
git push origin "v${VERSION}"

echo ""
echo -e "${GREEN}✅ Release v${VERSION} triggered successfully!${NC}"
echo ""
echo -e "${BLUE}📋 What happens next:${NC}"
echo "1. 🔨 GitHub Actions will build the binaries"
echo "2. 📦 Create release packages (.tar.gz files)"
echo "3.  Create GitHub release with download links"
echo "4. 🍺 Generate Homebrew formula automatically"
echo ""
echo -e "${BLUE}🔗 Monitor the release:${NC}"
echo "https://github.com/rosaboyle/leanmcp-cli-chat-deploy/actions"
echo "https://github.com/rosaboyle/leanmcp-cli-chat-deploy/releases"
echo ""
echo -e "${YELLOW}⏱️  Release should be ready in 5-10 minutes${NC}"
echo ""
echo -e "${GREEN}🎉 Users can install with:${NC}"
echo "brew tap rosaboyle/leanmcp && brew install leanmcp"
