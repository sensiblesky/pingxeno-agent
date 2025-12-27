package agent

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
)

// Identity represents the agent's system identity
type Identity struct {
	Hostname  string
	OSType    string
	OSVersion string
	IPAddress string
	MachineID string // Linux systemd machine-id or equivalent
	SystemUUID string // BIOS/UEFI system UUID
	DiskUUID  string // Root filesystem/disk UUID
	AgentID   string // Persistent agent ID (hash of machine-id + system_uuid + disk_uuid)
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

	// Get persistent identifiers
	machineID := getMachineID()
	systemUUID := getSystemUUID()
	diskUUID := getDiskUUID()
	
	// Generate persistent agent ID from identifiers
	agentID := generateAgentID(machineID, systemUUID, diskUUID)

	return &Identity{
		Hostname:   hostname,
		OSType:     osType,
		OSVersion:  osVersion,
		IPAddress:  ipAddress,
		MachineID:  machineID,
		SystemUUID: systemUUID,
		DiskUUID:   diskUUID,
		AgentID:    agentID,
	}, nil
}

// getMachineID retrieves the machine ID (Linux systemd, macOS, Windows, etc.)
func getMachineID() string {
	switch runtime.GOOS {
	case "linux":
		// Try systemd machine-id first
		if data, err := ioutil.ReadFile("/etc/machine-id"); err == nil {
			return strings.TrimSpace(string(data))
		}
		// Fallback to dbus machine-id
		if data, err := ioutil.ReadFile("/var/lib/dbus/machine-id"); err == nil {
			return strings.TrimSpace(string(data))
		}
	case "darwin":
		// macOS: Use system serial number as machine ID
		cmd := exec.Command("system_profiler", "SPHardwareDataType")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Serial Number") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	case "windows":
		// Windows: Use machine GUID from registry
		cmd := exec.Command("powershell", "-Command", "Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Cryptography' -Name MachineGuid | Select-Object -ExpandProperty MachineGuid")
		if output, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	case "freebsd":
		// FreeBSD: Use hostid
		cmd := exec.Command("hostid")
		if output, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	}
	return ""
}

// getSystemUUID retrieves the system UUID (BIOS/UEFI)
func getSystemUUID() string {
	switch runtime.GOOS {
	case "linux":
		// Try DMI system UUID
		if data, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
			uuid := strings.TrimSpace(string(data))
			if uuid != "" && uuid != "Not Specified" {
				return uuid
			}
		}
		// Alternative path
		if data, err := ioutil.ReadFile("/sys/class/dmi/id/board_id"); err == nil {
			uuid := strings.TrimSpace(string(data))
			if uuid != "" && uuid != "Not Specified" {
				return uuid
			}
		}
		// Try dmidecode as fallback
		cmd := exec.Command("dmidecode", "-s", "system-uuid")
		if output, err := cmd.Output(); err == nil {
			uuid := strings.TrimSpace(string(output))
			if uuid != "" && !strings.Contains(uuid, "Not Specified") {
				return uuid
			}
		}
	case "darwin":
		// macOS: Use system UUID
		cmd := exec.Command("system_profiler", "SPHardwareDataType")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Hardware UUID") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	case "windows":
		// Windows: Use system UUID from WMI
		cmd := exec.Command("powershell", "-Command", "Get-WmiObject Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID")
		if output, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	case "freebsd":
		// FreeBSD: Try DMI if available
		if data, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
			uuid := strings.TrimSpace(string(data))
			if uuid != "" && uuid != "Not Specified" {
				return uuid
			}
		}
	}
	return ""
}

// getDiskUUID retrieves the root filesystem/disk UUID
func getDiskUUID() string {
	switch runtime.GOOS {
	case "linux":
		// Try to get root filesystem UUID
		// Method 1: From /etc/fstab or findmnt
		cmd := exec.Command("findmnt", "-n", "-o", "UUID", "/")
		if output, err := cmd.Output(); err == nil {
			uuid := strings.TrimSpace(string(output))
			if uuid != "" {
				return uuid
			}
		}
		// Method 2: From blkid
		cmd = exec.Command("blkid", "-s", "UUID", "-o", "value", "-t", "TYPE!=swap")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			if len(lines) > 0 && lines[0] != "" {
				return lines[0]
			}
		}
		// Method 3: From /dev/disk/by-uuid (get first non-swap disk)
		files, err := ioutil.ReadDir("/dev/disk/by-uuid")
		if err == nil {
			for _, file := range files {
				uuid := file.Name()
				if uuid != "" {
					return uuid
				}
			}
		}
	case "darwin":
		// macOS: Get disk UUID
		cmd := exec.Command("diskutil", "info", "/")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Volume UUID") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	case "windows":
		// Windows: Get disk serial number
		cmd := exec.Command("powershell", "-Command", "Get-Disk | Where-Object {$_.IsBoot -eq $true} | Select-Object -ExpandProperty SerialNumber")
		if output, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	case "freebsd":
		// FreeBSD: Get disk UUID from geom
		cmd := exec.Command("geom", "disk", "list")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "ident:") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	}
	return ""
}

// generateAgentID creates a persistent agent ID from machine identifiers
func generateAgentID(machineID, systemUUID, diskUUID string) string {
	// Combine all identifiers
	combined := fmt.Sprintf("%s|%s|%s", machineID, systemUUID, diskUUID)
	
	// Generate SHA256 hash
	hash := sha256.Sum256([]byte(combined))
	
	// Return first 32 characters of hex representation
	return fmt.Sprintf("%x", hash)[:32]
}

