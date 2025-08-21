# LeanMCP CLI (`leanmcp`)

A command-line interface for interacting with LeanMCP APIs. Manage your projects, chats, deployments, and API keys from the terminal.

## 🚀 Installation

### Build from Source

```bash
git clone <repository-url>
cd leanmcp-cli-chat-deploy
go mod tidy
make build
```

### Create Alias (Optional)

Add this to your shell profile (`.bashrc`, `.zshrc`, etc.):

```bash
alias lcli='leanmcp'
```

## 🔐 Authentication

Before using the CLI, you need to authenticate with your API key:

```bash
leanmcp auth login --api-key airtrain_your_key_here
```

Your API key will be securely stored in `~/.leanmcp/config.yaml`.

### Authentication Commands

```bash
# Login with API key
leanmcp auth login --api-key <your-key>

# Check authentication status
leanmcp auth whoami

# Test API connection
leanmcp auth status

# Logout (remove stored credentials)
leanmcp auth logout
```

## 📋 Project Management

```bash
# List all projects
leanmcp projects list

# Show project details
leanmcp projects show <project-id>

# Create a new project
leanmcp projects create --name "My Project" --description "Optional description"

# Delete a project (requires --force)
leanmcp projects delete <project-id> --force

# List project builds
leanmcp projects builds <project-id>

# Start a new build
leanmcp projects build <project-id>
```

## 💬 Chat Management

```bash
# List all chats
leanmcp chats list

# Show chat details
leanmcp chats show <chat-id>

# Show chat message history
leanmcp chats history <chat-id>

# Limit message history
leanmcp chats history <chat-id> --limit 10

# Create a new chat
leanmcp chats create --title "My Chat" --model "gpt-4"

# Delete a chat (requires --force)
leanmcp chats delete <chat-id> --force
```

## 🔑 API Key Management

```bash
# List current API key info
leanmcp api-keys list

# Show detailed API key information
leanmcp api-keys info
```

## 🚀 Deployments

```bash
# List deployments (coming soon)
leanmcp deployments list

# Show deployment details (coming soon)
leanmcp deployments show <deployment-id>

# Show deployment logs (coming soon)
leanmcp deployments logs <deployment-id>
```

## ⚙️ Configuration

The CLI stores configuration in `~/.leanmcp/config.yaml`:

```yaml
api_key: "encrypted_key_here"
user_email: "user@example.com"
scopes: "BUILD_AND_DEPLOY,CHAT"
stored_at: "2025-01-09T21:40:36Z"
base_url: "https://api.leanmcp.ai"
```

### Global Flags

```bash
--config string     config file (default is $HOME/.leanmcp/config.yaml)
--base-url string   API base URL (default: https://api.leanmcp.ai)
--verbose, -v       verbose output
```

## 🎨 Output Formats

The CLI provides clean, colorized output with:

- ✅ **Success messages** in green
- ⚠️ **Warnings** in yellow  
- ❌ **Errors** in red
- 📋 **Tables** for list commands
- 🔍 **Detailed views** for specific items

## 🔧 Development

### Project Structure

```
leanmcp/
├── cmd/                 # Cobra commands
│   ├── root.go         # Root command
│   ├── auth.go         # Authentication commands
│   ├── projects.go     # Project commands
│   ├── chats.go        # Chat commands
│   ├── api-keys.go     # API key commands
│   └── deployments.go  # Deployment commands
├── internal/
│   ├── api/            # API client
│   ├── auth/           # Authentication management
│   ├── config/         # Configuration management
│   └── display/        # Output formatting
├── main.go             # Entry point
├── go.mod
└── README.md
```

### Build Commands

```bash
# Development build
go build -o leanmcp .

# Run tests
go test ./...

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o leanmcp-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o leanmcp-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o leanmcp-windows-amd64.exe .
```

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For support and questions, please contact the development team or create an issue in the repository.
