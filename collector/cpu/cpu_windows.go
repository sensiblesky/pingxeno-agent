//go:build windows

package cpu

import (
	"github.com/shirou/gopsutil/cpu"
)

type WindowsCollector struct{}

func newCollector() Collector {
	return &WindowsCollector{}
}

func (c *WindowsCollector) GetUsagePercent() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, nil
}

func (c *WindowsCollector) GetCores() (int, error) {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *WindowsCollector) GetLoadAvg() (load1, load5, load15 float64, err error) {
	// Windows doesn't have load average in the same way
	// Return 0 for all values
	return 0, 0, 0, nil
}

