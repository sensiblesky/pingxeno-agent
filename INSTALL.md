# Installation Guide

## Installing Go

### macOS (using Homebrew)
```bash
brew install go
```

### macOS (manual installation)
1. Download Go from https://go.dev/dl/
2. Install the package
3. Verify installation:
   ```bash
   go version
   ```

### Linux
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang-go

# Or download from https://go.dev/dl/
```

### Windows
1. Download the installer from https://go.dev/dl/
2. Run the installer
3. Add Go to your PATH if not done automatically

## Setting Up the Agent

Once Go is installed:

```bash
cd agent

# Download dependencies
go mod tidy

# Verify the build
go build ./cmd/agent

# Or run directly
go run ./cmd/agent/main.go --help
```

## Quick Test

After installing Go and running `go mod tidy`, test the agent:

```bash
# Test that it compiles
go build ./cmd/agent

# Run the help command
./agent --help

# Or if on Windows
./agent.exe --help
```

## Troubleshooting

If you get import errors:
1. Make sure you're in the `agent` directory
2. Run `go mod tidy` to download dependencies
3. If gopsutil errors persist, try:
   ```bash
   go get github.com/shirou/gopsutil@v2.21.11+incompatible
   go mod tidy
   ```

