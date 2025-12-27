package protocol

import "time"

// MetricsPayload represents the data sent to the API
type MetricsPayload struct {
	ServerKey            string                 `json:"server_key"`
	CPUUsagePercent      *float64               `json:"cpu_usage_percent,omitempty"`
	CPUCores             *int                   `json:"cpu_cores,omitempty"`
	CPULoad1Min          *float64               `json:"cpu_load_1min,omitempty"`
	CPULoad5Min          *float64               `json:"cpu_load_5min,omitempty"`
	CPULoad15Min         *float64               `json:"cpu_load_15min,omitempty"`
	MemoryTotalBytes     *int64                 `json:"memory_total_bytes,omitempty"`
	MemoryUsedBytes      *int64                 `json:"memory_used_bytes,omitempty"`
	MemoryFreeBytes      *int64                 `json:"memory_free_bytes,omitempty"`
	MemoryUsagePercent   *float64               `json:"memory_usage_percent,omitempty"`
	SwapTotalBytes       *int64                 `json:"swap_total_bytes,omitempty"`
	SwapUsedBytes        *int64                 `json:"swap_used_bytes,omitempty"`
	SwapFreeBytes        *int64                 `json:"swap_free_bytes,omitempty"`
	SwapUsagePercent     *float64               `json:"swap_usage_percent,omitempty"`
	DiskUsage            []DiskPartition        `json:"disk_usage,omitempty"`
	DiskTotalBytes       *int64                 `json:"disk_total_bytes,omitempty"`
	DiskUsedBytes        *int64                 `json:"disk_used_bytes,omitempty"`
	DiskFreeBytes        *int64                 `json:"disk_free_bytes,omitempty"`
	DiskUsagePercent     *float64               `json:"disk_usage_percent,omitempty"`
	NetworkInterfaces    []NetworkInterface     `json:"network_interfaces,omitempty"`
	NetworkBytesSent     *int64                 `json:"network_bytes_sent,omitempty"`
	NetworkBytesReceived *int64                 `json:"network_bytes_received,omitempty"`
	NetworkPacketsSent   *int64                 `json:"network_packets_sent,omitempty"`
	NetworkPacketsReceived *int64              `json:"network_packets_received,omitempty"`
	UptimeSeconds        *int                   `json:"uptime_seconds,omitempty"`
	ProcessesTotal       *int                   `json:"processes_total,omitempty"`
	ProcessesRunning     *int                   `json:"processes_running,omitempty"`
	ProcessesSleeping     *int                   `json:"processes_sleeping,omitempty"`
	Processes            []Process               `json:"processes,omitempty"`
	Hostname             string                 `json:"hostname,omitempty"`
	OSType               string                 `json:"os_type,omitempty"`
	OSVersion            string                 `json:"os_version,omitempty"`
	IPAddress            string                 `json:"ip_address,omitempty"`
	AgentVersion         string                 `json:"agent_version,omitempty"`
	RecordedAt           time.Time              `json:"recorded_at"`
}

// DiskPartition represents disk partition information
type DiskPartition struct {
	Device     string  `json:"device"`
	MountPoint string  `json:"mount_point"`
	FSType     string  `json:"fs_type"`
	TotalBytes int64   `json:"total_bytes"`
	UsedBytes  int64   `json:"used_bytes"`
	FreeBytes  int64   `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

// NetworkInterface represents network interface statistics
type NetworkInterface struct {
	Name            string `json:"name"`
	BytesSent       int64  `json:"bytes_sent"`
	BytesReceived   int64  `json:"bytes_received"`
	PacketsSent     int64  `json:"packets_sent"`
	PacketsReceived int64  `json:"packets_received"`
	ErrorsIn        int64  `json:"errors_in"`
	ErrorsOut       int64  `json:"errors_out"`
	DropIn          int64  `json:"drop_in"`
	DropOut         int64  `json:"drop_out"`
}

// Process represents process information
type Process struct {
	PID         int     `json:"pid"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryBytes int64   `json:"memory_bytes"`
	User        string  `json:"user"`
	Command     string  `json:"command"`
	CreatedAt   int64   `json:"created_at"`
}

