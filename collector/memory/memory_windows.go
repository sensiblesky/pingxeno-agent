//go:build windows

package memory

import (
	"github.com/shirou/gopsutil/mem"
)

type WindowsCollector struct{}

func newCollector() Collector {
	return &WindowsCollector{}
}

func (c *WindowsCollector) GetMemory() (total, used, free uint64, err error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return v.Total, v.Used, v.Available, nil
}

func (c *WindowsCollector) GetSwap() (total, used, free uint64, err error) {
	s, err := mem.SwapMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return s.Total, s.Used, s.Free, nil
}

