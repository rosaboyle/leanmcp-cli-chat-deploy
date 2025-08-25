# LeanMCP CLI Makefile

.PHONY: build clean install test help run

# Default target
all: build

# Build the CLI
build:
	@echo "🔨 Building LeanMCP CLI..."
	go build -o leanmcp-cli .
	@echo "✅ Build complete!"

# Build with custom version
build-version:
	@echo "🔨 Building LeanMCP CLI with version $(VERSION)..."
	@if [ -z "$(VERSION)" ]; then echo "Error: VERSION not set. Use: make build-version VERSION=1.0.0"; exit 1; fi
	go build -ldflags "-X 'github.com/ddod/leanmcp-cli/cmd.Version=$(VERSION)'" -o leanmcp-cli .
	@echo "✅ Build complete! Version: $(VERSION)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f leanmcp-cli
	rm -f leanmcp-cli-*
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

# Cross-platform builds
build-all: clean
	@echo "🌍 Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build -o leanmcp-cli-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o leanmcp-cli-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o leanmcp-cli-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o leanmcp-cli-windows-amd64.exe .
	@echo "✅ Cross-platform builds complete!"

# Development run
run: build
	@echo " Running CLI..."
	./leanmcp-cli

# Show help
help:
	@echo "LeanMCP CLI Development Commands:"
	@echo ""
	@echo "  build      - Build the CLI binary"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install Go dependencies"
	@echo "  test       - Run tests"
	@echo "  build-all  - Build for all platforms"
	@echo "  run        - Build and run CLI"
	@echo "  help       - Show this help"
	@echo ""
	@echo "Example usage:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make build-all"
