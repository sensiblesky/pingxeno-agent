package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/pingxeno/agent/collector/cpu"
	"github.com/pingxeno/agent/collector/disk"
	"github.com/pingxeno/agent/collector/memory"
	"github.com/pingxeno/agent/collector/network"
	"github.com/pingxeno/agent/collector/process"
	"github.com/pingxeno/agent/config"
	"github.com/pingxeno/agent/protocol"
	"github.com/pingxeno/agent/scheduler"
	"github.com/pingxeno/agent/sender"
	"github.com/shirou/gopsutil/host"
	"go.uber.org/zap"
)

// Agent represents the monitoring agent
type Agent struct {
	config     *config.Config
	scheduler  *scheduler.Scheduler
	sender     *sender.RetrySender
	identity   *Identity
	logger     *zap.Logger
	cpuCol     cpu.Collector
	memCol     memory.Collector
	diskCol    disk.Collector
	netCol     network.Collector
	procCol    process.Collector
}

// NewAgent creates a new agent instance
func NewAgent(cfg *config.Config, logger *zap.Logger) (*Agent, error) {
	identity, err := GetIdentity()
	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}

	httpClient := sender.NewClient(cfg, logger)
	retrySender := sender.NewRetrySender(
		httpClient,
		cfg.Sender.RetryAttempts,
		cfg.Sender.RetryBackoff,
		logger,
	)

	sch := scheduler.NewScheduler(cfg.Collection.Interval, cfg.Collection.Jitter)

	agent := &Agent{
		config:    cfg,
		scheduler: sch,
		sender:    retrySender,
		identity:  identity,
		logger:    logger,
		cpuCol:    cpu.NewCollector(),
		memCol:    memory.NewCollector(),
		diskCol:   disk.NewCollector(),
		netCol:    network.NewCollector(),
		procCol:   process.NewCollector(),
	}

	return agent, nil
}

// CollectMetrics collects all system metrics
func (a *Agent) CollectMetrics() (*protocol.MetricsPayload, error) {
	payload := &protocol.MetricsPayload{
		ServerKey:  a.config.Server.ServerKey,
		Hostname:   a.identity.Hostname,
		OSType:     a.identity.OSType,
		OSVersion:  a.identity.OSVersion,
		IPAddress:  a.identity.IPAddress,
		RecordedAt: time.Now(),
	}

	// Collect CPU metrics
	cpuMetrics, err := cpu.Collect(a.cpuCol)
	if err != nil {
		a.logger.Warn("Failed to collect CPU metrics", zap.Error(err))
	} else {
		payload.CPUUsagePercent = &cpuMetrics.UsagePercent
		payload.CPUCores = &cpuMetrics.Cores
		payload.CPULoad1Min = &cpuMetrics.Load1Min
		payload.CPULoad5Min = &cpuMetrics.Load5Min
		payload.CPULoad15Min = &cpuMetrics.Load15Min
	}

	// Collect Memory metrics
	memMetrics, err := memory.Collect(a.memCol)
	if err != nil {
		a.logger.Warn("Failed to collect memory metrics", zap.Error(err))
	} else {
		payload.MemoryTotalBytes = &memMetrics.MemoryTotalBytes
		payload.MemoryUsedBytes = &memMetrics.MemoryUsedBytes
		payload.MemoryFreeBytes = &memMetrics.MemoryFreeBytes
		payload.MemoryUsagePercent = &memMetrics.MemoryUsagePercent
		payload.SwapTotalBytes = &memMetrics.SwapTotalBytes
		payload.SwapUsedBytes = &memMetrics.SwapUsedBytes
		payload.SwapFreeBytes = &memMetrics.SwapFreeBytes
		payload.SwapUsagePercent = &memMetrics.SwapUsagePercent
	}

	// Collect Disk metrics
	diskMetrics, err := disk.Collect(a.diskCol)
	if err != nil {
		a.logger.Warn("Failed to collect disk metrics", zap.Error(err))
	} else {
		payload.DiskUsage = diskMetrics.Partitions
		payload.DiskTotalBytes = &diskMetrics.TotalBytes
		payload.DiskUsedBytes = &diskMetrics.UsedBytes
		payload.DiskFreeBytes = &diskMetrics.FreeBytes
		payload.DiskUsagePercent = &diskMetrics.UsagePercent
	}

	// Collect Network metrics
	netMetrics, err := network.Collect(a.netCol)
	if err != nil {
		a.logger.Warn("Failed to collect network metrics", zap.Error(err))
	} else {
		payload.NetworkInterfaces = netMetrics.Interfaces
		payload.NetworkBytesSent = &netMetrics.BytesSent
		payload.NetworkBytesReceived = &netMetrics.BytesReceived
		payload.NetworkPacketsSent = &netMetrics.PacketsSent
		payload.NetworkPacketsReceived = &netMetrics.PacketsReceived
	}

	// Collect Process metrics
	procMetrics, err := process.Collect(a.procCol)
	if err != nil {
		a.logger.Warn("Failed to collect process metrics", zap.Error(err))
	} else {
		payload.ProcessesTotal = &procMetrics.Total
		payload.ProcessesRunning = &procMetrics.Running
		payload.ProcessesSleeping = &procMetrics.Sleeping
		payload.Processes = procMetrics.Processes
	}

	// Get uptime
	uptime, err := host.Uptime()
	if err == nil {
		uptimeInt := int(uptime)
		payload.UptimeSeconds = &uptimeInt
	}

	return payload, nil
}

// Run starts the agent's main loop
func (a *Agent) Run(ctx context.Context) error {
	a.logger.Info("Agent started",
		zap.String("hostname", a.identity.Hostname),
		zap.String("os", a.identity.OSType),
		zap.String("api_url", a.config.Server.APIURL),
	)

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Agent stopping")
			return nil
		default:
			// Collect metrics
			payload, err := a.CollectMetrics()
			if err != nil {
				a.logger.Error("Failed to collect metrics", zap.Error(err))
				a.scheduler.Wait()
				continue
			}

			// Send metrics
			if err := a.sender.SendWithRetry(payload); err != nil {
				a.logger.Error("Failed to send metrics", zap.Error(err))
			} else {
				a.logger.Debug("Metrics sent successfully")
			}

			// Wait for next collection
			a.scheduler.Wait()
		}
	}
}

// TestConnection tests the connection to the API
func (a *Agent) TestConnection() error {
	client := sender.NewClient(a.config, a.logger)
	return client.TestConnection()
}

// Sender returns the agent's sender instance (for testing)
func (a *Agent) Sender() *sender.RetrySender {
	return a.sender
}

