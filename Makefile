# LeanMCP CLI Makefile
# Build targets for Homebrew distribution

.PHONY: build clean install test help run release

# Version (can be overridden)
VERSION ?= 1.0.0

# Default target
all: build

# Detect OS-specific binary name
BINARY_NAME = leanmcp
ifeq ($(OS),Windows_NT)
    BINARY_NAME := leanmcp.exe
endif

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -ldflags "-X github.com/ddod/leanmcp-cli/cmd.Version=$(VERSION)" -o $(BINARY_NAME) .
	@echo "✅ Build complete! Version: $(VERSION)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f leanmcp
	rm -f leanmcp-*
	rm -f *.tar.gz
	@echo "✅ Clean complete!"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed!"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...
	@echo "✅ Tests complete!"

# Build for Homebrew (macOS architectures)
build-homebrew: clean
	@echo "🍺 Building for Homebrew (macOS only)..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'github.com/ddod/leanmcp-cli/cmd.Version=$(VERSION)'" -o leanmcp-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'github.com/ddod/leanmcp-cli/cmd.Version=$(VERSION)'" -o leanmcp-darwin-arm64 .
	@echo "✅ Homebrew builds complete!"

# Create release packages
release: build-homebrew
	@echo "📦 Creating release packages..."
	tar -czf leanmcp-$(VERSION)-darwin-amd64.tar.gz leanmcp-darwin-amd64
	tar -czf leanmcp-$(VERSION)-darwin-arm64.tar.gz leanmcp-darwin-arm64
	@echo "✅ Release packages created:"
	@ls -la leanmcp-$(VERSION)-*.tar.gz

# Development run
run: build
	@echo " Running leanmcp..."
	./leanmcp

# Show help
help:
	@echo "LeanMCP CLI Development Commands:"
	@echo ""
	@echo "  build           - Build leanmcp binary"
	@echo "  clean           - Clean build artifacts"
	@echo "  install         - Install Go dependencies"
	@echo "  test            - Run tests"
	@echo "  build-homebrew  - Build for macOS (Intel + Apple Silicon)"
	@echo "  release         - Create release packages for Homebrew"
	@echo "  run             - Build and run leanmcp"
	@echo "  help            - Show this help"
	@echo ""
	@echo "Example usage:"
	@echo "  make build VERSION=1.0.0"
	@echo "  make release VERSION=1.0.0"
