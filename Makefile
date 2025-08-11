# MCLI - LeanMCP CLI Makefile
# Build targets for Homebrew distribution

.PHONY: build clean install test help run release

# Version (can be overridden)
VERSION ?= 1.0.0

# Default target
all: build

# Build the CLI as mcli
build:
	@echo "🔨 Building mcli..."
	go build -ldflags "-X 'github.com/ddod/leanmcp-cli/cmd.Version=$(VERSION)'" -o mcli .
	@echo "✅ Build complete! Version: $(VERSION)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f mcli
	rm -f mcli-*
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
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'github.com/ddod/leanmcp-cli/cmd.Version=$(VERSION)'" -o mcli-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'github.com/ddod/leanmcp-cli/cmd.Version=$(VERSION)'" -o mcli-darwin-arm64 .
	@echo "✅ Homebrew builds complete!"

# Create release packages
release: build-homebrew
	@echo "📦 Creating release packages..."
	tar -czf mcli-$(VERSION)-darwin-amd64.tar.gz mcli-darwin-amd64
	tar -czf mcli-$(VERSION)-darwin-arm64.tar.gz mcli-darwin-arm64
	@echo "✅ Release packages created:"
	@ls -la mcli-$(VERSION)-*.tar.gz

# Development run
run: build
	@echo "🚀 Running mcli..."
	./mcli

# Show help
help:
	@echo "MCLI Development Commands:"
	@echo ""
	@echo "  build           - Build mcli binary"
	@echo "  clean           - Clean build artifacts"
	@echo "  install         - Install Go dependencies"
	@echo "  test            - Run tests"
	@echo "  build-homebrew  - Build for macOS (Intel + Apple Silicon)"
	@echo "  release         - Create release packages for Homebrew"
	@echo "  run             - Build and run mcli"
	@echo "  help            - Show this help"
	@echo ""
	@echo "Example usage:"
	@echo "  make build VERSION=1.0.0"
	@echo "  make release VERSION=1.0.0"
