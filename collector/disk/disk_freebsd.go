//go:build freebsd

package disk

import (
	"github.com/pingxeno/agent/protocol"
	"github.com/shirou/gopsutil/disk"
)

type FreeBSDCollector struct{}

func newCollector() Collector {
	return &FreeBSDCollector{}
}

func (c *FreeBSDCollector) GetPartitions() ([]protocol.DiskPartition, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var result []protocol.DiskPartition
	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}

		result = append(result, protocol.DiskPartition{
			Device:      part.Device,
			MountPoint:  part.Mountpoint,
			FSType:      part.Fstype,
			TotalBytes:  int64(usage.Total),
			UsedBytes:   int64(usage.Used),
			FreeBytes:   int64(usage.Free),
			UsagePercent: usage.UsedPercent,
		})
	}

	return result, nil
}

func (c *FreeBSDCollector) GetUsage(path string) (total, used, free uint64, err error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return 0, 0, 0, err
	}
	return usage.Total, usage.Used, usage.Free, nil
}

