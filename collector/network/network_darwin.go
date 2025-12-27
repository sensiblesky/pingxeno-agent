//go:build darwin

package network

import (
	"github.com/pingxeno/agent/protocol"
	"github.com/shirou/gopsutil/net"
)

type DarwinCollector struct{}

func newCollector() Collector {
	return &DarwinCollector{}
}

func (c *DarwinCollector) GetInterfaces() ([]protocol.NetworkInterface, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var result []protocol.NetworkInterface
	for _, stat := range stats {
		result = append(result, protocol.NetworkInterface{
			Name:            stat.Name,
			BytesSent:       int64(stat.BytesSent),
			BytesReceived:   int64(stat.BytesRecv),
			PacketsSent:     int64(stat.PacketsSent),
			PacketsReceived: int64(stat.PacketsRecv),
			ErrorsIn:        int64(stat.Errin),
			ErrorsOut:       int64(stat.Errout),
			DropIn:          int64(stat.Dropin),
			DropOut:         int64(stat.Dropout),
		})
	}

	return result, nil
}

func (c *DarwinCollector) GetIO() (bytesSent, bytesRecv, packetsSent, packetsRecv int64, err error) {
	stats, err := net.IOCounters(false)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if len(stats) > 0 {
		stat := stats[0]
		return int64(stat.BytesSent), int64(stat.BytesRecv), int64(stat.PacketsSent), int64(stat.PacketsRecv), nil
	}

	return 0, 0, 0, 0, nil
}

