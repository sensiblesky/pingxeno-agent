//go:build freebsd

package cpu

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

type FreeBSDCollector struct{}

func newCollector() Collector {
	return &FreeBSDCollector{}
}

func (c *FreeBSDCollector) GetUsagePercent() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, nil
}

func (c *FreeBSDCollector) GetCores() (int, error) {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *FreeBSDCollector) GetLoadAvg() (load1, load5, load15 float64, err error) {
	avg, err := load.Avg()
	if err != nil {
		return 0, 0, 0, err
	}
	return avg.Load1, avg.Load5, avg.Load15, nil
}

