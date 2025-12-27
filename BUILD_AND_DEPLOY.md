# Building and Deploying the Agent Binary

## Building on Dev Server (Recommended)

You can build the binary on your development machine and transfer only the compiled binary to the remote server. **No Go or dependencies are needed on the remote server.**

### Step 1: Build the Binary

On your **development machine** (where Go is installed):

```bash
cd /path/to/pingxeno-agent

# Build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o pingxeno-agent-linux-amd64 ./cmd/agent

# Or build for Linux ARM64
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o pingxeno-agent-linux-arm64 ./cmd/agent

# Make executable
chmod +x pingxeno-agent-linux-amd64
```

**Flags explained:**
- `-ldflags="-s -w"` - Strips debug symbols to reduce binary size (optional but recommended)
- `GOOS=linux` - Target operating system
- `GOARCH=amd64` or `arm64` - Target architecture

### Step 2: Transfer Binary to Remote Server

```bash
# Using SCP
scp pingxeno-agent-linux-amd64 user@remote-server:/tmp/pingxeno-agent

# Or using rsync (more reliable)
rsync -avz pingxeno-agent-linux-amd64 user@remote-server:/tmp/pingxeno-agent
```

### Step 3: Install on Remote Server

**On the remote server** (no Go needed):

```bash
# Move to system location
sudo mv /tmp/pingxeno-agent /usr/local/bin/pingxeno-agent
sudo chmod +x /usr/local/bin/pingxeno-agent

# Verify it works
pingxeno-agent --help
```

That's it! The binary is completely standalone and requires **nothing** on the remote server.

## Why This Works

Go compiles everything into a **static binary**:
- ✅ All dependencies are included
- ✅ No runtime dependencies needed
- ✅ No Go installation required
- ✅ No shared libraries needed (unless using CGO)

## Binary Size

The compiled binary is typically:
- **Without flags**: ~15-20 MB
- **With `-ldflags="-s -w"`**: ~10-15 MB (smaller, no debug info)

## Verification

To verify the binary is standalone:

```bash
# Check if it's a static binary
file pingxeno-agent-linux-amd64
# Should show: ELF 64-bit LSB executable, x86-64, statically linked

# Check dependencies (should show minimal or none)
ldd pingxeno-agent-linux-amd64
# Should show: not a dynamic executable (for static builds)
```

## Complete Deployment Example

### On Dev Machine:

```bash
# 1. Build
cd pingxeno-agent
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o pingxeno-agent-linux-amd64 ./cmd/agent
chmod +x pingxeno-agent-linux-amd64

# 2. Transfer
scp pingxeno-agent-linux-amd64 root@192.168.1.100:/tmp/
```

### On Remote Server:

```bash
# 1. Install binary
sudo mv /tmp/pingxeno-agent-linux-amd64 /usr/local/bin/pingxeno-agent
sudo chmod +x /usr/local/bin/pingxeno-agent

# 2. Create config
sudo mkdir -p /etc/pingxeno
sudo nano /etc/pingxeno/agent.yaml
# (Add your configuration)

# 3. Create systemd service
sudo nano /etc/systemd/system/pingxeno-agent.service
# (Add service configuration)

# 4. Start service
sudo systemctl daemon-reload
sudo systemctl start pingxeno-agent
sudo systemctl enable pingxeno-agent
```

## Advantages of Building on Dev Server

✅ **Faster deployment** - No need to install Go on remote server  
✅ **Smaller footprint** - Only transfer the binary (~10-15 MB)  
✅ **No dependencies** - Remote server doesn't need Go, Git, or any build tools  
✅ **Consistent builds** - Build once, deploy everywhere  
✅ **Security** - No need to expose build tools on production servers  

## When You DO Need Dependencies

You only need Go and dependencies on the remote server if:
- ❌ You want to build from source on the remote server
- ❌ You're using CGO with dynamic linking (not the case here)
- ❌ You want to modify and rebuild on the remote server

For normal deployment, **just transfer the binary** - it's completely self-contained!

