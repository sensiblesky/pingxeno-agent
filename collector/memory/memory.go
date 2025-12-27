package memory

// Collector interface for memory metrics
type Collector interface {
	GetMemory() (total, used, free uint64, err error)
	GetSwap() (total, used, free uint64, err error)
}

// Metrics represents memory metrics
type Metrics struct {
	MemoryTotalBytes   int64
	MemoryUsedBytes    int64
	MemoryFreeBytes    int64
	MemoryUsagePercent float64
	SwapTotalBytes     int64
	SwapUsedBytes      int64
	SwapFreeBytes      int64
	SwapUsagePercent   float64
}

// NewCollector creates a platform-specific memory collector
func NewCollector() Collector {
	return newCollector()
}

// Collect gathers all memory metrics
func Collect(c Collector) (*Metrics, error) {
	memTotal, memUsed, memFree, err := c.GetMemory()
	if err != nil {
		return nil, err
	}

	swapTotal, swapUsed, swapFree, err := c.GetSwap()
	if err != nil {
		// Swap might not be available on all systems
		swapTotal, swapUsed, swapFree = 0, 0, 0
	}

	memTotalInt := int64(memTotal)
	memUsedInt := int64(memUsed)
	memFreeInt := int64(memFree)
	swapTotalInt := int64(swapTotal)
	swapUsedInt := int64(swapUsed)
	swapFreeInt := int64(swapFree)

	var memUsagePercent float64
	if memTotalInt > 0 {
		memUsagePercent = (float64(memUsedInt) / float64(memTotalInt)) * 100
	}

	var swapUsagePercent float64
	if swapTotalInt > 0 {
		swapUsagePercent = (float64(swapUsedInt) / float64(swapTotalInt)) * 100
	}

	return &Metrics{
		MemoryTotalBytes:   memTotalInt,
		MemoryUsedBytes:    memUsedInt,
		MemoryFreeBytes:    memFreeInt,
		MemoryUsagePercent: memUsagePercent,
		SwapTotalBytes:     swapTotalInt,
		SwapUsedBytes:      swapUsedInt,
		SwapFreeBytes:      swapFreeInt,
		SwapUsagePercent:   swapUsagePercent,
	}, nil
}
