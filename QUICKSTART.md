# Quick Start Guide

## Step 1: Install Go

### macOS
```bash
# Using Homebrew (recommended)
brew install go

# Or download from https://go.dev/dl/
```

### Verify Go Installation
```bash
go version
# Should output: go version go1.21.x ...
```

## Step 2: Setup the Agent

```bash
# Navigate to agent directory
cd agent

# Download all dependencies
go mod tidy

# This will download:
# - github.com/shirou/gopsutil (system metrics)
# - github.com/spf13/cobra (CLI framework)
# - github.com/spf13/viper (configuration)
# - go.uber.org/zap (logging)
```

## Step 3: Build the Agent

```bash
# Build for your current platform
go build -o pingxeno-agent ./cmd/agent

# Or build for specific platforms
GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux ./cmd/agent
GOOS=windows GOARCH=amd64 go build -o pingxeno-agent-windows.exe ./cmd/agent
GOOS=darwin GOARCH=amd64 go build -o pingxeno-agent-macos ./cmd/agent
```

## Step 4: Configure the Agent

```bash
# Interactive installation (prompts for API URL, keys, etc.)
./pingxeno-agent install

# Or manually create config file
cat > ~/.config/pingxeno/agent.yaml << EOF
server:
  api_url: "https://your-domain.com/api/v1/server-stats"
  api_key: "pk_your_api_key_here"
  server_key: "srv_your_server_key_here"

collection:
  interval: 60s
  jitter: 5s

sender:
  batch_size: 10
  retry_attempts: 3
  retry_backoff: 2s

security:
  tls_skip_verify: false
  timeout: 30s

logging:
  level: "info"
  file: ""
EOF
```

## Step 5: Test the Agent

```bash
# Check status and test connection
./pingxeno-agent status

# Send a test metrics payload
./pingxeno-agent test

# Run the agent (foreground)
./pingxeno-agent run
```

## Step 6: Run as Service (Optional)

### Linux (systemd)
```bash
sudo ./pingxeno-agent install --config /etc/pingxeno/agent.yaml
sudo systemctl start pingxeno-agent
sudo systemctl enable pingxeno-agent
```

### macOS (LaunchAgent)
```bash
./pingxeno-agent install --config ~/Library/LaunchAgents/pingxeno-agent.yaml
launchctl load ~/Library/LaunchAgents/pingxeno-agent.plist
```

### Windows (Service)
```powershell
.\pingxeno-agent.exe install --config "C:\Program Files\PingXeno\agent.yaml"
```

## Troubleshooting

### "command not found: go"
- Install Go first (see Step 1)
- Make sure Go is in your PATH: `echo $PATH | grep go`

### Import errors after `go mod tidy`
```bash
# Clear module cache and re-download
go clean -modcache
go mod tidy
```

### gopsutil import errors
```bash
# Force download specific version
go get github.com/shirou/gopsutil@v2.21.11+incompatible
go mod tidy
```

### Build errors
```bash
# Check Go version (needs 1.21+)
go version

# Update Go if needed
brew upgrade go  # macOS
```

## Next Steps

- Read [README.md](README.md) for full documentation
- See [INSTALL.md](INSTALL.md) for detailed installation instructions
- Check [IMPORT_FIX.md](IMPORT_FIX.md) if you encounter import issues

