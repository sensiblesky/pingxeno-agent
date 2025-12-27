# Installing PingXeno Agent from GitHub Repository

This guide will help you clone the repository and install the agent on a remote server.

## Prerequisites

- Go 1.21 or higher installed on the server
- Git installed
- Access to the GitHub repository
- API key and server key from your PingXeno dashboard

## Step 1: Clone the Repository

```bash
# Clone the repository
git clone https://github.com/sensiblesky/pingxeno-agent.git
cd pingxeno-agent

# Or if using SSH
git clone git@github.com:sensiblesky/pingxeno-agent.git
cd pingxeno-agent
```

## Step 2: Install Go (if not already installed)

### Linux (Ubuntu/Debian)
```bash
# Remove old Go if exists
sudo rm -rf /usr/local/go

# Download and install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

### Linux (CentOS/RHEL)
```bash
# Download and install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bash_profile
source ~/.bash_profile

# Verify installation
go version
```

## Step 3: Build the Agent

```bash
# Navigate to agent directory
cd pingxeno-agent

# Download dependencies
go mod download

# Build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux-amd64 ./cmd/agent

# Or build for Linux ARM64 (for Raspberry Pi, ARM servers)
GOOS=linux GOARCH=arm64 go build -o pingxeno-agent-linux-arm64 ./cmd/agent

# Make executable
chmod +x pingxeno-agent-linux-amd64
```

## Step 4: Get Server Credentials

1. Log in to your PingXeno dashboard
2. Navigate to **Monitoring** → **Server Monitoring**
3. Create a new server or select an existing one
4. Copy the following:
   - **Server Key**: `srv_xxxxxxxxxxxxx`
   - **API Key**: `pk_xxxxxxxxxxxxx` (from the linked API key)
   - **API Endpoint**: `https://your-domain.com/api/v1/server-stats`

## Step 5: Install and Configure

### Option A: Interactive Installation

```bash
# Run the agent with install command (if implemented)
./pingxeno-agent-linux-amd64 install

# Or configure manually
./pingxeno-agent-linux-amd64 config
```

### Option B: Manual Configuration

```bash
# Create configuration directory
sudo mkdir -p /etc/pingxeno

# Create configuration file
sudo nano /etc/pingxeno/agent.yaml
```

Add the following content:

```yaml
server:
  api_url: "https://your-domain.com/api/v1/server-stats"
  api_key: "pk_your_api_key_here"
  server_key: "srv_your_server_key_here"
  
collection:
  interval: 60s
  jitter: 5s
  
sender:
  batch_size: 10
  batch_timeout: 30s
  retry_attempts: 3
  retry_backoff: 2s
  
security:
  tls_skip_verify: false
  timeout: 30s
  
logging:
  level: "info"
  file: "/var/log/pingxeno-agent.log"
```

Save and exit (Ctrl+X, then Y, then Enter).

## Step 6: Move Binary to System Location

```bash
# Move binary to system location
sudo mv pingxeno-agent-linux-amd64 /usr/local/bin/pingxeno-agent
sudo chmod +x /usr/local/bin/pingxeno-agent

# Verify it works
pingxeno-agent --help
```

## Step 7: Create Systemd Service

```bash
# Create systemd service file
sudo nano /etc/systemd/system/pingxeno-agent.service
```

Add the following content:

```ini
[Unit]
Description=PingXeno Monitoring Agent
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/pingxeno-agent run --config /etc/pingxeno/agent.yaml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Save and exit.

## Step 8: Start and Enable the Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Start the service
sudo systemctl start pingxeno-agent

# Enable to start on boot
sudo systemctl enable pingxeno-agent

# Check status
sudo systemctl status pingxeno-agent

# View logs
sudo journalctl -u pingxeno-agent -f
```

## Step 9: Verify Installation

```bash
# Check if agent is running
sudo systemctl status pingxeno-agent

# Check logs for any errors
sudo journalctl -u pingxeno-agent -n 50

# Test connection (if implemented)
pingxeno-agent test
```

## Troubleshooting

### Agent Not Starting

```bash
# Check service status
sudo systemctl status pingxeno-agent

# Check logs
sudo journalctl -u pingxeno-agent -n 100

# Check configuration file
cat /etc/pingxeno/agent.yaml

# Test binary manually
sudo /usr/local/bin/pingxeno-agent run --config /etc/pingxeno/agent.yaml
```

### Connection Issues

```bash
# Test API endpoint connectivity
curl -v https://your-domain.com/api/v1/server-stats

# Check firewall
sudo ufw status
sudo iptables -L

# Verify API key and server key are correct
cat /etc/pingxeno/agent.yaml | grep -E "api_key|server_key"
```

### Build Errors

```bash
# Clean and rebuild
go clean -cache
go mod tidy
go build -o pingxeno-agent-linux-amd64 ./cmd/agent

# Check Go version
go version  # Should be 1.21 or higher
```

## Updating the Agent

```bash
# Navigate to repository
cd ~/pingxeno-agent

# Pull latest changes
git pull origin main

# Rebuild
go mod download
GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux-amd64 ./cmd/agent

# Stop service
sudo systemctl stop pingxeno-agent

# Replace binary
sudo mv pingxeno-agent-linux-amd64 /usr/local/bin/pingxeno-agent
sudo chmod +x /usr/local/bin/pingxeno-agent

# Restart service
sudo systemctl start pingxeno-agent

# Check status
sudo systemctl status pingxeno-agent
```

## Uninstalling

```bash
# Stop and disable service
sudo systemctl stop pingxeno-agent
sudo systemctl disable pingxeno-agent

# Remove service file
sudo rm /etc/systemd/system/pingxeno-agent.service
sudo systemctl daemon-reload

# Remove binary
sudo rm /usr/local/bin/pingxeno-agent

# Remove configuration (optional)
sudo rm -rf /etc/pingxeno

# Remove logs (optional)
sudo rm /var/log/pingxeno-agent.log
```

## Quick Install Script

For convenience, you can use this quick install script:

```bash
#!/bin/bash
set -e

echo "Installing PingXeno Agent..."

# Clone repository
if [ ! -d "pingxeno-agent" ]; then
    git clone https://github.com/sensiblesky/pingxeno-agent.git
fi
cd pingxeno-agent

# Build
echo "Building agent..."
go mod download
GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux-amd64 ./cmd/agent
chmod +x pingxeno-agent-linux-amd64

# Install
echo "Installing agent..."
sudo mv pingxeno-agent-linux-amd64 /usr/local/bin/pingxeno-agent

# Create config directory
sudo mkdir -p /etc/pingxeno

echo "Installation complete!"
echo "Next steps:"
echo "1. Configure: sudo nano /etc/pingxeno/agent.yaml"
echo "2. Create systemd service: sudo nano /etc/systemd/system/pingxeno-agent.service"
echo "3. Start: sudo systemctl start pingxeno-agent"
```

Save as `install.sh`, make executable (`chmod +x install.sh`), and run (`./install.sh`).

