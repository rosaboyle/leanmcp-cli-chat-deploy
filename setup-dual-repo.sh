#!/bin/bash

# Dual Repository Setup for MCLI
# Keeps source private, makes releases public

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🏗️ Setting up Dual Repository Architecture${NC}"
echo "=============================================="

echo -e "${YELLOW}📋 Setup Instructions:${NC}"
echo ""

echo -e "${BLUE}Step 1: Create Public Release Repository${NC}"
echo "1. Go to: https://github.com/new"
echo "2. Repository name: ${BLUE}mcli-releases${NC}"
echo "3. Description: \"Public releases for MCLI\""
echo "4. Make it ${BLUE}PUBLIC${NC}"
echo "5. Initialize with README"
echo ""

echo -e "${BLUE}Step 2: Generate Personal Access Token${NC}"
echo "1. Go to: https://github.com/settings/tokens"
echo "2. Generate new token (classic)"
echo "3. Scopes needed: ${BLUE}public_repo, workflow${NC}"
echo "4. Copy the token"
echo ""

echo -e "${BLUE}Step 3: Add Secret to Private Repository${NC}"
echo "1. Go to: https://github.com/rosaboyle/leanmcp-cli-chat-deploy/settings/secrets/actions"
echo "2. Add secret: ${BLUE}RELEASE_TOKEN${NC}"
echo "3. Value: Your personal access token from Step 2"
echo ""

echo -e "${BLUE}Files to be created:${NC}"
echo "- .github/workflows/dual-repo-release.yml (Updated CI/CD)"
echo "- dual-repo-release.sh (Release script for dual repos)"
echo "- Updated Homebrew formula pointing to public repo"

echo ""
echo -e "${GREEN}🎯 Result: Source stays private, releases are public!${NC}"
echo "- Private repo: rosaboyle/leanmcp-cli-chat-deploy"
echo "- Public repo:  rosaboyle/mcli-releases"
echo "- Homebrew downloads from public repo only"
