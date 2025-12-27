# PingXeno Agent

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS%20%7C%20FreeBSD-lightgrey)](https://github.com/sensiblesky/pingxeno-agent)

A lightweight, cross-platform system monitoring agent written in Go that collects comprehensive server performance metrics and sends them to the [PingXeno](https://github.com/sensiblesky/pingxeno) monitoring platform.

## 🚀 Features

- **📊 Comprehensive Metrics Collection**
  - CPU usage, load averages, and core count
  - Memory and swap utilization
  - Disk usage per partition with filesystem details
  - Network interface statistics (bytes, packets, errors)
  - Process monitoring with detailed information (PID, CPU%, Memory%, user, command)
  - System uptime tracking

- **🌐 Cross-Platform Support**
  - ✅ Linux (AMD64, ARM64)
  - ✅ Windows (AMD64)
  - ✅ macOS (AMD64, ARM64 - Apple Silicon)
  - ✅ FreeBSD

- **🔒 Security & Reliability**
  - API key authentication
  - TLS/SSL encryption support
  - Offline buffering with retry logic
  - Batch sending for efficiency
  - Configurable timeouts and retry policies

- **⚙️ Easy Configuration**
  - Interactive setup wizard
  - YAML configuration file
  - Environment variable support
  - Command-line flags

- **🔄 Service Integration**
  - Systemd service (Linux)
  - Windows Service
  - LaunchAgent (macOS)
  - Manual/foreground mode

- **📈 Production Ready**
  - Structured logging with rotation
  - Health checks and diagnostics
  - Graceful shutdown handling
  - Resource-efficient design

## 📋 Prerequisites

- **Go 1.21 or higher** (for building from source)
- Network access to your PingXeno server
- API key and server key from your PingXeno dashboard

## 🏃 Quick Start

### Download Pre-built Binary

Visit the [Releases](https://github.com/sensiblesky/pingxeno-agent/releases) page to download the binary for your platform.

### Install and Configure

```bash
# Make executable (Linux/macOS)
chmod +x pingxeno-agent

# Interactive configuration
./pingxeno-agent install

# Or configure manually
./pingxeno-agent config
```

### Run the Agent

```bash
# Run in foreground (for testing)
./pingxeno-agent run

# Install as system service (Linux)
sudo ./pingxeno-agent install --service

# Check status
./pingxeno-agent status
```

## 📖 Installation Guide

### Linux

```bash
# Download latest release
curl -L https://github.com/sensiblesky/pingxeno-agent/releases/latest/download/pingxeno-agent-linux-amd64 -o /usr/local/bin/pingxeno-agent
chmod +x /usr/local/bin/pingxeno-agent

# Install and configure
sudo pingxeno-agent install

# Start service
sudo systemctl start pingxeno-agent
sudo systemctl enable pingxeno-agent
```

### Windows

```powershell
# Download latest release
Invoke-WebRequest -Uri "https://github.com/sensiblesky/pingxeno-agent/releases/latest/download/pingxeno-agent-windows-amd64.exe" -OutFile "C:\Program Files\PingXeno\pingxeno-agent.exe"

# Install and configure
pingxeno-agent.exe install

# Start service
Start-Service PingXenoAgent
```

### macOS

```bash
# Download latest release
curl -L https://github.com/sensiblesky/pingxeno-agent/releases/latest/download/pingxeno-agent-darwin-amd64 -o /usr/local/bin/pingxeno-agent
chmod +x /usr/local/bin/pingxeno-agent

# Install and configure
pingxeno-agent install

# Start service (LaunchAgent)
launchctl load ~/Library/LaunchAgents/com.pingxeno.agent.plist
```

## ⚙️ Configuration

### Configuration File (`agent.yaml`)

```yaml
server:
  api_url: "https://your-domain.com/api/v1/server-stats"
  api_key: "pk_your_api_key_here"
  server_key: "srv_your_server_key_here"
  
collection:
  interval: 60s          # How often to collect metrics
  jitter: 5s             # Random jitter to prevent thundering herd
  
sender:
  batch_size: 10         # Number of metrics to batch before sending
  batch_timeout: 30s     # Max time to wait before sending batch
  retry_attempts: 3      # Number of retry attempts on failure
  retry_backoff: 2s      # Initial backoff time between retries
  
security:
  tls_skip_verify: false # Skip TLS certificate verification (not recommended)
  timeout: 30s           # HTTP request timeout
  
logging:
  level: "info"          # Log level: debug, info, warn, error
  file: "/var/log/pingxeno-agent.log"  # Log file path
```

### Environment Variables

```bash
export PINGXENO_API_URL="https://your-domain.com/api/v1/server-stats"
export PINGXENO_API_KEY="pk_your_api_key_here"
export PINGXENO_SERVER_KEY="srv_your_server_key_here"
export PINGXENO_COLLECTION_INTERVAL="60s"
```

### Command-Line Flags

```bash
./pingxeno-agent run \
  --api-url "https://your-domain.com/api/v1/server-stats" \
  --api-key "pk_your_api_key_here" \
  --server-key "srv_your_server_key_here" \
  --interval "60s"
```

## 📊 Metrics Collected

### CPU Metrics
- CPU usage percentage
- CPU cores count
- Load averages (1min, 5min, 15min) - Linux/Unix only

### Memory Metrics
- Total, used, and free memory (bytes)
- Memory usage percentage
- Swap statistics (total, used, free, percentage)

### Disk Metrics
- Per-partition details:
  - Device name
  - Mount point
  - Filesystem type
  - Total, used, free space (bytes)
  - Usage percentage
- Overall disk statistics

### Network Metrics
- Per-interface statistics:
  - Bytes sent/received
  - Packets sent/received
  - Errors and drops (in/out)
- Aggregate network statistics

### Process Metrics
- Total, running, and sleeping processes
- Detailed process list:
  - PID, name, status
  - CPU and memory usage percentages
  - Memory usage in bytes
  - User and command

### System Metrics
- System uptime (seconds)
- Hostname, OS type, OS version
- IP address
- Agent version

## 🛠️ Building from Source

### Prerequisites

- Go 1.21 or higher
- Git

### Build Steps

```bash
# Clone the repository
git clone https://github.com/sensiblesky/pingxeno-agent.git
cd pingxeno-agent

# Download dependencies
go mod download

# Build for current platform
go build -o pingxeno-agent ./cmd/agent

# Cross-compile for specific platform
GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux-amd64 ./cmd/agent
GOOS=windows GOARCH=amd64 go build -o pingxeno-agent-windows-amd64.exe ./cmd/agent
GOOS=darwin GOARCH=amd64 go build -o pingxeno-agent-darwin-amd64 ./cmd/agent
GOOS=darwin GOARCH=arm64 go build -o pingxeno-agent-darwin-arm64 ./cmd/agent

# Build all platforms (if build script exists)
./scripts/build.sh --all
```

## 📝 Usage

### Commands

```bash
# Run agent in foreground
./pingxeno-agent run

# Install and configure interactively
./pingxeno-agent install

# Configure settings
./pingxeno-agent config

# Check agent status
./pingxeno-agent status

# View logs
./pingxeno-agent logs

# Test connection
./pingxeno-agent test

# Uninstall service
./pingxeno-agent uninstall

# Show version
./pingxeno-agent version
```

### Service Management (Linux)

```bash
# Start service
sudo systemctl start pingxeno-agent

# Stop service
sudo systemctl stop pingxeno-agent

# Restart service
sudo systemctl restart pingxeno-agent

# Check status
sudo systemctl status pingxeno-agent

# View logs
sudo journalctl -u pingxeno-agent -f
```

## 🔍 Troubleshooting

### Common Issues

**Agent not sending data:**
- Verify API URL is correct and accessible
- Check API key and server key are valid
- Ensure firewall allows outbound HTTPS connections
- Check agent logs: `./pingxeno-agent logs` or `journalctl -u pingxeno-agent`

**Connection errors:**
- Verify TLS certificate is valid (or set `tls_skip_verify: true` for testing)
- Check network connectivity: `curl -v https://your-domain.com/api/v1/server-stats`
- Verify API key has correct scopes (needs `create` scope)

**High CPU/Memory usage:**
- Increase collection interval (e.g., `interval: 300s` for 5 minutes)
- Reduce process list size in configuration
- Check for too many network interfaces or disk partitions

### Debug Mode

```bash
# Run with debug logging
./pingxeno-agent run --log-level debug

# Test connection without sending data
./pingxeno-agent test
```

### Logs Location

- **Linux (systemd)**: `journalctl -u pingxeno-agent`
- **Linux (manual)**: `/var/log/pingxeno-agent.log` (or configured path)
- **Windows**: `C:\ProgramData\PingXeno\logs\agent.log`
- **macOS**: `~/Library/Logs/pingxeno-agent.log`

## 🔐 Security Considerations

1. **API Keys**: Store API keys securely. Use environment variables or secure configuration files with restricted permissions (e.g., `chmod 600 agent.yaml`)

2. **TLS**: Always use HTTPS in production. Only disable TLS verification (`tls_skip_verify: true`) for testing

3. **File Permissions**: Ensure configuration files are readable only by the agent user:
   ```bash
   chmod 600 /etc/pingxeno/agent.yaml
   chown pingxeno:pingxeno /etc/pingxeno/agent.yaml
   ```

4. **Network**: The agent only makes outbound HTTPS connections. No inbound ports are required

## 📚 Documentation

- [Installation Guide](INSTALL.md)
- [Linux Deployment Guide](DEPLOY_LINUX.md)
- [Troubleshooting Guide](TROUBLESHOOTING.md)
- [Quick Start Guide](QUICKSTART.md)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [PingXeno Platform](https://github.com/sensiblesky/pingxeno) - The monitoring platform that receives metrics from this agent

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/sensiblesky/pingxeno-agent/issues)
- **Documentation**: [Wiki](https://github.com/sensiblesky/pingxeno-agent/wiki)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [gopsutil](https://github.com/shirou/gopsutil) for cross-platform system metrics
- Inspired by modern monitoring agents like Datadog, New Relic, and Prometheus exporters

---

**Made with ❤️ by the PingXeno team**
