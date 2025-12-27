#!/bin/bash
set -e

echo "========================================="
echo "PingXeno Agent - Quick Install Script"
echo "========================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "Please install Go 1.21 or higher first."
    echo "Visit: https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Check if Git is installed
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed!"
    echo "Please install Git first."
    exit 1
fi

echo "✅ Git is installed"
echo ""

# Clone or update repository
if [ -d "pingxeno-agent" ]; then
    echo "📦 Updating existing repository..."
    cd pingxeno-agent
    git pull origin main
else
    echo "📦 Cloning repository..."
    git clone https://github.com/sensiblesky/pingxeno-agent.git
    cd pingxeno-agent
fi

# Download dependencies
echo "📥 Downloading dependencies..."
go mod download

# Build agent
echo "🔨 Building agent..."
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    GOOS=linux GOARCH=amd64 go build -o pingxeno-agent-linux-amd64 ./cmd/agent
    BINARY="pingxeno-agent-linux-amd64"
elif [ "$ARCH" = "aarch64" ]; then
    GOOS=linux GOARCH=arm64 go build -o pingxeno-agent-linux-arm64 ./cmd/agent
    BINARY="pingxeno-agent-linux-arm64"
else
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
fi

chmod +x $BINARY

# Install binary
echo "📦 Installing binary..."
sudo mv $BINARY /usr/local/bin/pingxeno-agent
sudo chmod +x /usr/local/bin/pingxeno-agent

echo ""
echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "1. Configure the agent:"
echo "   sudo nano /etc/pingxeno/agent.yaml"
echo ""
echo "2. Create systemd service:"
echo "   sudo nano /etc/systemd/system/pingxeno-agent.service"
echo ""
echo "3. Start the service:"
echo "   sudo systemctl daemon-reload"
echo "   sudo systemctl start pingxeno-agent"
echo "   sudo systemctl enable pingxeno-agent"
echo ""
echo "4. Check status:"
echo "   sudo systemctl status pingxeno-agent"
echo ""
