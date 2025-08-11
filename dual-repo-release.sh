#!/bin/bash

# Dual Repository Release Script for MCLI
# Releases from private repo to public repo automatically

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check if version is provided
if [ "$#" -ne 1 ]; then
    echo -e "${RED}❌ Usage: $0 <version>${NC}"
    echo -e "${BLUE}Example: $0 1.0.2${NC}"
    exit 1
fi

VERSION=$1

# Validate version format
if [[ ! $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Invalid version format. Use semantic versioning (e.g., 1.0.2)${NC}"
    exit 1
fi

echo -e "${BLUE}🚀 Dual Repository Release v${VERSION}${NC}"
echo "=============================="

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}❌ Not a git repository${NC}"
    exit 1
fi

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo -e "${YELLOW}⚠️  Uncommitted changes detected${NC}"
    echo -e "${BLUE}📝 Committing all changes...${NC}"
    git add .
    git commit -m "Prepare release v${VERSION}"
fi

# Update version in source code
echo -e "${BLUE}📝 Updating version in source code...${NC}"
sed -i.bak "s/Version = \".*\"/Version = \"${VERSION}\"/" cmd/root.go
rm -f cmd/root.go.bak

# Commit version update
git add cmd/root.go
git commit -m "Bump version to ${VERSION}" || echo "No version changes to commit"

# Create and push tag
echo -e "${BLUE}🏷️  Creating git tag v${VERSION}...${NC}"
git tag "v${VERSION}"

echo -e "${BLUE}⬆️  Pushing to private repository...${NC}"
git push origin main
git push origin "v${VERSION}"

echo -e "${GREEN}✅ Release initiated!${NC}"
echo ""
echo -e "${BLUE}🔄 What happens next (automatically):${NC}"
echo "1. 🏗️  GitHub Actions builds binaries"
echo "2. 📦 Creates release in PUBLIC repository: rosaboyle/mcli-releases"
echo "3. 🍺 Updates Homebrew formula automatically"
echo "4. ✅ Users can install: brew tap rosaboyle/mcli && brew install mcli"
echo ""
echo -e "${YELLOW}📊 Monitor progress:${NC}"
echo "- Private repo actions: https://github.com/rosaboyle/leanmcp-cli-chat-deploy/actions"
echo "- Public releases: https://github.com/rosaboyle/mcli-releases/releases"
echo "- Homebrew formula: https://github.com/rosaboyle/homebrew-mcli/blob/main/mcli.rb"
echo ""
echo -e "${GREEN}🎉 Release v${VERSION} in progress!${NC}"
