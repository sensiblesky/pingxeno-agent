//go:build darwin

package cpu

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

type DarwinCollector struct{}

func newCollector() Collector {
	return &DarwinCollector{}
}

func (c *DarwinCollector) GetUsagePercent() (float64, error) {
	// On macOS, cpu.Percent needs a non-zero interval to calculate usage
	// Use 1 second interval for accurate measurement
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, nil
}

func (c *DarwinCollector) GetCores() (int, error) {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *DarwinCollector) GetLoadAvg() (load1, load5, load15 float64, err error) {
	avg, err := load.Avg()
	if err != nil {
		return 0, 0, 0, err
	}
	return avg.Load1, avg.Load5, avg.Load15, nil
}

