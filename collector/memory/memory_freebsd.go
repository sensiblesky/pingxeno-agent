//go:build freebsd

package memory

import (
	"github.com/shirou/gopsutil/mem"
)

type FreeBSDCollector struct{}

func newCollector() Collector {
	return &FreeBSDCollector{}
}

func (c *FreeBSDCollector) GetMemory() (total, used, free uint64, err error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return v.Total, v.Used, v.Available, nil
}

func (c *FreeBSDCollector) GetSwap() (total, used, free uint64, err error) {
	s, err := mem.SwapMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return s.Total, s.Used, s.Free, nil
}

