package cpu

// Collector interface for CPU metrics
type Collector interface {
	GetUsagePercent() (float64, error)
	GetCores() (int, error)
	GetLoadAvg() (load1, load5, load15 float64, err error)
}

// Metrics represents CPU metrics
type Metrics struct {
	UsagePercent float64
	Cores        int
	Load1Min     float64
	Load5Min     float64
	Load15Min    float64
}

// NewCollector creates a platform-specific CPU collector
func NewCollector() Collector {
	return newCollector()
}

// Collect gathers all CPU metrics
func Collect(c Collector) (*Metrics, error) {
	usage, err := c.GetUsagePercent()
	if err != nil {
		return nil, err
	}

	cores, err := c.GetCores()
	if err != nil {
		return nil, err
	}

	load1, load5, load15, err := c.GetLoadAvg()
	if err != nil {
		// Load average might not be available on all systems
		load1, load5, load15 = 0, 0, 0
	}

	return &Metrics{
		UsagePercent: usage,
		Cores:        cores,
		Load1Min:     load1,
		Load5Min:     load5,
		Load15Min:    load15,
	}, nil
}
