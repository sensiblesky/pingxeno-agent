//go:build !linux && !windows && !darwin && !freebsd

package process

import (
	"github.com/pingxeno/agent/protocol"
	"github.com/shirou/gopsutil/process"
)

type DefaultCollector struct{}

func newCollector() Collector {
	return &DefaultCollector{}
}

func (c *DefaultCollector) GetProcessCount() (total, running, sleeping int, err error) {
	pids, err := process.Pids()
	if err != nil {
		return 0, 0, 0, err
	}

	total = len(pids)
	running = 0
	sleeping = 0

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		status, err := p.Status()
		if err != nil {
			continue
		}

		switch status[0] {
		case 'R': // Running
			running++
		case 'S', 'D': // Sleeping, Disk sleep
			sleeping++
		}
	}

	return total, running, sleeping, nil
}

func (c *DefaultCollector) GetAllProcesses() ([]protocol.Process, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	var processes []protocol.Process
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		name, _ := p.Name()
		status, _ := p.Status()
		cpuPercent, _ := p.CPUPercent()
		memPercent32, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		username, _ := p.Username()
		cmdline, _ := p.Cmdline()
		createTime, _ := p.CreateTime()

		statusChar := ""
		if len(status) > 0 {
			statusChar = string(status[0])
		}

		var memBytes int64
		if memInfo != nil {
			memBytes = int64(memInfo.RSS)
		}

		// Convert float32 to float64
		memPercent := float64(memPercent32)

		processes = append(processes, protocol.Process{
			PID:          int(pid),
			Name:         name,
			Status:       statusChar,
			CPUPercent:   cpuPercent,
			MemoryPercent: memPercent,
			MemoryBytes:  memBytes,
			User:         username,
			Command:      cmdline,
			CreatedAt:    createTime,
		})
	}

	return processes, nil
}

