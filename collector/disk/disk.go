package disk

import (
	"github.com/pingxeno/agent/protocol"
)

// Collector interface for disk metrics
type Collector interface {
	GetPartitions() ([]protocol.DiskPartition, error)
	GetUsage(path string) (total, used, free uint64, err error)
}

// Metrics represents disk metrics
type Metrics struct {
	Partitions    []protocol.DiskPartition
	TotalBytes    int64
	UsedBytes     int64
	FreeBytes     int64
	UsagePercent  float64
}

// NewCollector creates a platform-specific disk collector
func NewCollector() Collector {
	return newCollector()
}

// Collect gathers all disk metrics
func Collect(c Collector) (*Metrics, error) {
	partitions, err := c.GetPartitions()
	if err != nil {
		return nil, err
	}

	// Virtual filesystem mount points to exclude (common across all OS)
	virtualMountPoints := map[string]bool{
		"/dev":      true,
		"/proc":     true,
		"/sys":     true,
		"/run":     true,
		"/tmp":     true,
		"/snap":    true,
		"/boot/efi": true,
		"/var/run": true,
		"/var/lock": true,
		"/boot":    true, // Boot partition is separate
		"/swap":    true,
	}
	
	// Virtual filesystem types to exclude
	virtualFSTypes := map[string]bool{
		"devtmpfs":   true,
		"devfs":      true,
		"proc":       true,
		"procfs":     true,
		"sysfs":      true,
		"tmpfs":      true,
		"overlay":    true,
		"cgroup":     true,
		"cgroup2":   true,
		"pstore":     true,
		"bpf":        true,
		"tracefs":    true,
		"debugfs":    true,
		"securityfs": true,
		"hugetlbfs":  true,
		"mqueue":     true,
		"systemd-1":  true,
		"binfmt_misc": true,
		"fusectl":    true,
		"configfs":   true,
		"autofs":     true,
		"rpc_pipefs": true,
		"nfsd":       true,
		"none":       true,
		"swap":       true,
	}
	
	// Find root filesystem (/) - this is what we'll use for totals
	var rootPartition *protocol.DiskPartition
	var allPartitions []protocol.DiskPartition
	
	for _, part := range partitions {
		// Always include in partitions list for display
		allPartitions = append(allPartitions, part)
		
		// Skip virtual filesystems
		if virtualMountPoints[part.MountPoint] || virtualFSTypes[part.FSType] {
			continue
		}
		
		// Skip if mount point contains virtual filesystem paths
		skip := false
		for vmp := range virtualMountPoints {
			if len(part.MountPoint) >= len(vmp) && part.MountPoint[:len(vmp)] == vmp {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		
		// Find root filesystem (/)
		if part.MountPoint == "/" {
			rootPartition = &part
			break // Root filesystem found, use this for totals
		}
	}
	
	// If root filesystem not found, try to find the largest physical partition
	if rootPartition == nil {
		var largestPartition *protocol.DiskPartition
		var largestSize int64 = 0
		
		for i := range allPartitions {
			part := &allPartitions[i]
			
			// Skip virtual filesystems
			if virtualMountPoints[part.MountPoint] || virtualFSTypes[part.FSType] {
				continue
			}
			
			// Skip if mount point contains virtual filesystem paths
			skip := false
			for vmp := range virtualMountPoints {
				if len(part.MountPoint) >= len(vmp) && part.MountPoint[:len(vmp)] == vmp {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
			
			// Find largest physical partition (likely the main disk)
			if part.TotalBytes > largestSize {
				largestSize = part.TotalBytes
				largestPartition = part
			}
		}
		
		if largestPartition != nil {
			rootPartition = largestPartition
		}
	}
	
	// Calculate totals from root filesystem only
	var totalBytes, usedBytes, freeBytes int64
	var usagePercent float64
	
	if rootPartition != nil {
		totalBytes = rootPartition.TotalBytes
		usedBytes = rootPartition.UsedBytes
		freeBytes = rootPartition.FreeBytes
		if totalBytes > 0 {
			usagePercent = (float64(usedBytes) / float64(totalBytes)) * 100
		}
	}

	return &Metrics{
		Partitions:   allPartitions, // Return all partitions for display
		TotalBytes:   totalBytes,    // But totals only from root filesystem
		UsedBytes:    usedBytes,
		FreeBytes:    freeBytes,
		UsagePercent: usagePercent,
	}, nil
}
