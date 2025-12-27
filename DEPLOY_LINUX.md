# Deploying PingXeno Agent on Linux Server

This guide will help you deploy the PingXeno agent on a remote Linux server.

## Prerequisites

1. **On your local machine (macOS/Windows):**
   - Go installed (for building the agent)
   - Access to the agent source code

2. **On the Linux server:**
   - SSH access
   - Basic Linux commands (curl, wget, etc.)

## Step 1: Build the Agent for Linux

On your local machine, navigate to the agent directory and build for Linux:

```bash
cd /Users/denicsann/Desktop/projects/pingxeno/agent

# Build for Linux AMD64 (most common)
GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux-amd64 ./cmd/agent

# Or build for Linux ARM64 (for Raspberry Pi, ARM servers)
GOOS=linux GOARCH=arm64 go build -o pingxeno-agent-linux-arm64 ./cmd/agent
```

The binary will be created in the current directory.

## Step 2: Get Server Credentials from Dashboard

1. Log in to your PingXeno dashboard: `http://127.0.0.1:8000` (or your domain)
2. Navigate to **Monitoring** → **Server Monitoring**
3. Create a new server or select an existing one
4. Copy the following from the server details page:
   - **Server Key**: `srv_xxxxxxxxxxxxx`
   - **API Key**: `pk_xxxxxxxxxxxxx` (from the linked API key)
   - **API Endpoint**: `http://127.0.0.1:8000/api/v1/server-stats` (or your domain)

## Step 3: Transfer the Binary to Linux Server

### Option A: Using SCP

```bash
# From your local machine
scp pingxeno-agent-linux-amd64 user@your-server-ip:/tmp/pingxeno-agent

# Example:
scp pingxeno-agent-linux-amd64 root@192.168.1.100:/tmp/pingxeno-agent
```

### Option B: Using wget/curl (if you host the binary)

```bash
# On the Linux server
wget https://your-domain.com/pingxeno-agent-linux-amd64 -O /tmp/pingxeno-agent
# or
curl -L https://your-domain.com/pingxeno-agent-linux-amd64 -o /tmp/pingxeno-agent
```

## Step 4: Install and Configure on Linux Server

SSH into your Linux server:

```bash
ssh user@your-server-ip
```

### Move the binary to a system location:

```bash
# Make it executable
chmod +x /tmp/pingxeno-agent

# Move to system location
sudo mv /tmp/pingxeno-agent /usr/local/bin/pingxeno-agent

# Verify it works
pingxeno-agent --help
```

### Configure the Agent

You have two options:

#### Option 1: Interactive Installation (Recommended)

```bash
sudo pingxeno-agent install
```

This will prompt you for:
- **API URL**: `http://127.0.0.1:8000/api/v1/server-stats` (or your domain)
- **API Key**: `pk_xxxxxxxxxxxxx`
- **Server Key**: `srv_xxxxxxxxxxxxx`
- **Collection Interval**: `60` (seconds, default)

The configuration will be saved to `/etc/pingxeno/agent.yaml` (or `~/.config/pingxeno/agent.yaml` if not root).

#### Option 2: Manual Configuration

Create the configuration file manually:

```bash
# Create config directory
sudo mkdir -p /etc/pingxeno

# Create config file
sudo nano /etc/pingxeno/agent.yaml
```

Add the following content (replace with your actual values):

```yaml
server:
  api_url: "http://127.0.0.1:8000/api/v1/server-stats"
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

## Step 5: Run the Agent

### Option 1: Run Manually (for testing)

```bash
# Run in foreground
sudo pingxeno-agent run --config /etc/pingxeno/agent.yaml

# Or run in background
sudo nohup pingxeno-agent run --config /etc/pingxeno/agent.yaml > /var/log/pingxeno-agent.log 2>&1 &
```

### Option 2: Install as Systemd Service (Recommended for Production)

Create a systemd service file:

```bash
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

Save and exit, then:

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable the service (start on boot)
sudo systemctl enable pingxeno-agent

# Start the service
sudo systemctl start pingxeno-agent

# Check status
sudo systemctl status pingxeno-agent

# View logs
sudo journalctl -u pingxeno-agent -f
```

## Step 6: Verify It's Working

1. **Check the service status:**
   ```bash
   sudo systemctl status pingxeno-agent
   ```

2. **Check the logs:**
   ```bash
   sudo journalctl -u pingxeno-agent -n 50
   # or
   tail -f /var/log/pingxeno-agent.log
   ```

3. **Test the connection:**
   ```bash
   sudo pingxeno-agent test --config /etc/pingxeno/agent.yaml
   ```

4. **Check the dashboard:**
   - Go to your PingXeno dashboard
   - Navigate to **Monitoring** → **Server Monitoring**
   - Click on your server
   - You should see metrics appearing within 1-2 minutes

## Troubleshooting

### Agent not sending data

1. **Check connectivity:**
   ```bash
   curl -X POST http://127.0.0.1:8000/api/v1/server-stats \
     -H "X-API-Key: pk_your_api_key" \
     -H "Content-Type: application/json" \
     -d '{"server_key":"srv_your_server_key","recorded_at":"2025-01-01T00:00:00Z"}'
   ```

2. **Check firewall:**
   ```bash
   # Allow outbound HTTPS/HTTP
   sudo ufw allow out 80/tcp
   sudo ufw allow out 443/tcp
   ```

3. **Check logs:**
   ```bash
   sudo journalctl -u pingxeno-agent -n 100 --no-pager
   ```

### Permission issues

```bash
# Make sure the binary is executable
sudo chmod +x /usr/local/bin/pingxeno-agent

# Check file permissions
ls -la /usr/local/bin/pingxeno-agent
ls -la /etc/pingxeno/agent.yaml
```

### Service won't start

```bash
# Check systemd logs
sudo journalctl -u pingxeno-agent -n 50

# Test manually
sudo /usr/local/bin/pingxeno-agent run --config /etc/pingxeno/agent.yaml
```

## Updating the Agent

1. Build a new binary on your local machine
2. Transfer it to the server
3. Stop the service:
   ```bash
   sudo systemctl stop pingxeno-agent
   ```
4. Replace the binary:
   ```bash
   sudo mv /tmp/pingxeno-agent-linux-amd64 /usr/local/bin/pingxeno-agent
   sudo chmod +x /usr/local/bin/pingxeno-agent
   ```
5. Start the service:
   ```bash
   sudo systemctl start pingxeno-agent
   ```

## Uninstalling

```bash
# Stop and disable the service
sudo systemctl stop pingxeno-agent
sudo systemctl disable pingxeno-agent

# Remove the service file
sudo rm /etc/systemd/system/pingxeno-agent.service
sudo systemctl daemon-reload

# Remove the binary
sudo rm /usr/local/bin/pingxeno-agent

# Remove configuration (optional)
sudo rm -rf /etc/pingxeno
```

## Quick Reference

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux-amd64 ./cmd/agent

# Transfer to server
scp pingxeno-agent-linux-amd64 user@server:/tmp/

# On server: Install
sudo mv /tmp/pingxeno-agent-linux-amd64 /usr/local/bin/pingxeno-agent
sudo chmod +x /usr/local/bin/pingxeno-agent
sudo pingxeno-agent install

# Run as service
sudo systemctl start pingxeno-agent
sudo systemctl status pingxeno-agent
```

## Notes

- Replace `http://127.0.0.1:8000` with your actual PingXeno domain if deployed
- For production, use HTTPS: `https://your-domain.com/api/v1/server-stats`
- The agent collects metrics every 60 seconds by default
- Metrics are sent in batches for efficiency
- The agent will retry failed requests automatically

