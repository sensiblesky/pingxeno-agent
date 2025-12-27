package agent

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
)

// Identity represents the agent's system identity
type Identity struct {
	Hostname  string
	OSType    string
	OSVersion string
	IPAddress string
}

// GetIdentity gathers system identity information
func GetIdentity() (*Identity, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	// Get primary IP address
	ipAddress := "unknown"
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			addrs := iface.Addrs
			if len(addrs) > 0 {
				for _, addr := range addrs {
					ip := addr.Addr
					if ip != "" && ip != "127.0.0.1" && ip != "::1" {
						ipAddress = ip
						break
					}
				}
				if ipAddress != "unknown" {
					break
				}
			}
		}
	}

	osType := runtime.GOOS
	osVersion := info.Platform + " " + info.PlatformVersion

	return &Identity{
		Hostname:  hostname,
		OSType:    osType,
		OSVersion: osVersion,
		IPAddress: ipAddress,
	}, nil
}

