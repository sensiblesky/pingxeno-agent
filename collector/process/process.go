package process

import "github.com/pingxeno/agent/protocol"

// Collector interface for process metrics
type Collector interface {
	GetProcessCount() (total, running, sleeping int, err error)
	GetAllProcesses() ([]protocol.Process, error)
}

// Metrics represents process metrics
type Metrics struct {
	Total    int
	Running  int
	Sleeping int
	Processes []protocol.Process
}

// NewCollector creates a platform-specific process collector
func NewCollector() Collector {
	return newCollector()
}

// Collect gathers all process metrics
func Collect(c Collector) (*Metrics, error) {
	total, running, sleeping, err := c.GetProcessCount()
	if err != nil {
		return nil, err
	}

	// Get all processes (limit to 1000 to avoid payload size issues)
	processes, err := c.GetAllProcesses()
	if err != nil {
		// Don't fail if we can't get process details, just log it
		processes = []protocol.Process{}
	}

	// Limit to 1000 processes to avoid payload size issues
	if len(processes) > 1000 {
		processes = processes[:1000]
	}

	return &Metrics{
		Total:     total,
		Running:   running,
		Sleeping:  sleeping,
		Processes: processes,
	}, nil
}
