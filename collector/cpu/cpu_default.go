//go:build !linux && !windows && !darwin && !freebsd

package cpu

import (
	"github.com/shirou/gopsutil/cpu"
)

type DefaultCollector struct{}

func newCollector() Collector {
	return &DefaultCollector{}
}

func (c *DefaultCollector) GetUsagePercent() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, nil
}

func (c *DefaultCollector) GetCores() (int, error) {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *DefaultCollector) GetLoadAvg() (load1, load5, load15 float64, err error) {
	// Not available on all platforms
	return 0, 0, 0, nil
}

