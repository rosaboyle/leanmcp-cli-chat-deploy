# LeanMCP CLI (`leanmcp-cli` / `lcli`)

A command-line interface for interacting with LeanMCP APIs. Manage your projects, chats, deployments, and API keys from the terminal.

## 🚀 Installation

### Build from Source

```bash
git clone <repository-url>
cd leanmcp-cli
go mod tidy
go build -o leanmcp-cli .
```

### Create Alias (Optional)

Add this to your shell profile (`.bashrc`, `.zshrc`, etc.):

```bash
alias lcli='leanmcp-cli'
```

## 🔐 Authentication

Before using the CLI, you need to authenticate with your API key:

```bash
leanmcp-cli auth login --api-key airtrain_your_key_here
```

Your API key will be securely stored in `~/.leanmcp-cli/config.yaml`.

### Authentication Commands

```bash
# Login with API key
leanmcp-cli auth login --api-key <your-key>

# Check authentication status
leanmcp-cli auth whoami

# Test API connection
leanmcp-cli auth status

# Logout (remove stored credentials)
leanmcp-cli auth logout
```

## 📋 Project Management

```bash
# List all projects
leanmcp-cli projects list

# Show project details
leanmcp-cli projects show <project-id>

# Create a new project
leanmcp-cli projects create --name "My Project" --description "Optional description"

# Delete a project (requires --force)
leanmcp-cli projects delete <project-id> --force

# List project builds
leanmcp-cli projects builds <project-id>

# Start a new build
leanmcp-cli projects build <project-id>
```

## 💬 Chat Management

```bash
# List all chats
leanmcp-cli chats list

# Show chat details
leanmcp-cli chats show <chat-id>

# Show chat message history
leanmcp-cli chats history <chat-id>

# Limit message history
leanmcp-cli chats history <chat-id> --limit 10

# Create a new chat
leanmcp-cli chats create --title "My Chat" --model "gpt-4"

# Delete a chat (requires --force)
leanmcp-cli chats delete <chat-id> --force
```

## 🔑 API Key Management

```bash
# List current API key info
leanmcp-cli api-keys list

# Show detailed API key information
leanmcp-cli api-keys info
```

## 🚀 Deployments

```bash
# List deployments (coming soon)
leanmcp-cli deployments list

# Show deployment details (coming soon)
leanmcp-cli deployments show <deployment-id>

# Show deployment logs (coming soon)
leanmcp-cli deployments logs <deployment-id>
```

## ⚙️ Configuration

The CLI stores configuration in `~/.leanmcp-cli/config.yaml`:

```yaml
api_key: "encrypted_key_here"
user_email: "user@example.com"
scopes: "BUILD_AND_DEPLOY,CHAT"
stored_at: "2025-01-09T21:40:36Z"
base_url: "https://api.leanmcp.ai"
```

### Global Flags

```bash
--config string     config file (default is $HOME/.leanmcp-cli/config.yaml)
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
leanmcp-cli/
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
go build -o leanmcp-cli .

# Run tests
go test ./...

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o leanmcp-cli-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o leanmcp-cli-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o leanmcp-cli-windows-amd64.exe .
```

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For support and questions, please contact the development team or create an issue in the repository.
