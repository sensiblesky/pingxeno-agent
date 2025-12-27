package network

import (
	"github.com/pingxeno/agent/protocol"
)

// Collector interface for network metrics
type Collector interface {
	GetInterfaces() ([]protocol.NetworkInterface, error)
	GetIO() (bytesSent, bytesRecv, packetsSent, packetsRecv int64, err error)
}

// Metrics represents network metrics
type Metrics struct {
	Interfaces         []protocol.NetworkInterface
	BytesSent          int64
	BytesReceived      int64
	PacketsSent         int64
	PacketsReceived     int64
}

// NewCollector creates a platform-specific network collector
func NewCollector() Collector {
	return newCollector()
}

// Collect gathers all network metrics
func Collect(c Collector) (*Metrics, error) {
	interfaces, err := c.GetInterfaces()
	if err != nil {
		return nil, err
	}

	bytesSent, bytesRecv, packetsSent, packetsRecv, err := c.GetIO()
	if err != nil {
		// IO stats might not be available
		bytesSent, bytesRecv, packetsSent, packetsRecv = 0, 0, 0, 0
	}

	return &Metrics{
		Interfaces:      interfaces,
		BytesSent:        bytesSent,
		BytesReceived:    bytesRecv,
		PacketsSent:      packetsSent,
		PacketsReceived:  packetsRecv,
	}, nil
}
