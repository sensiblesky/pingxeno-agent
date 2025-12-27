# Import Path Fixes

## Fixed Issues

1. **gopsutil imports**: Changed from `github.com/shirou/gopsutil/v3` to `github.com/shirou/gopsutil` (v2) for better compatibility
2. **Local module imports**: `github.com/pingxeno/agent/*` imports are correct - this is just the module name defined in `go.mod`

## About the Module Name

The module name `github.com/pingxeno/agent` doesn't need to exist on GitHub. It's just an identifier for your local Go module. All imports like:
- `github.com/pingxeno/agent/protocol`
- `github.com/pingxeno/agent/collector/cpu`
- etc.

Are relative to the module root defined in `go.mod`.

## Setup Instructions

1. Make sure you're in the `agent` directory
2. Run `go mod tidy` to download dependencies
3. If you still get errors, try:
   ```bash
   go get github.com/shirou/gopsutil@latest
   go get github.com/spf13/cobra@latest
   go get github.com/spf13/viper@latest
   go get go.uber.org/zap@latest
   go mod tidy
   ```

## Alternative Module Name

If you prefer a simpler module name, you can change the first line in `go.mod` to:
```
module pingxeno-agent
```

Then update all imports from `github.com/pingxeno/agent/...` to `pingxeno-agent/...`

