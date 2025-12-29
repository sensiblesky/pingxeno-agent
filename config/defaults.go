package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	var defaultLogFile string
	if runtime.GOOS == "windows" {
		// Windows: Use ProgramData directory
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = "C:\\ProgramData"
		}
		defaultLogFile = filepath.Join(programData, "PingXeno", "agent.log")
	} else {
		// Unix-like: Use /var/log
		defaultLogFile = "/var/log/pingxeno-agent.log"
	}

	return &Config{
		Server: ServerConfig{
			APIURL:    "http://localhost:8000/api/v1/server-stats",
			APIKey:    "",
			ServerKey: "",
		},
		Collection: CollectionConfig{
			Interval: 60 * time.Second,
			Jitter:   5 * time.Second,
		},
		Sender: SenderConfig{
			BatchSize:     10,
			BatchTimeout:  30 * time.Second,
			RetryAttempts: 3,
			RetryBackoff:  2 * time.Second,
		},
		Security: SecurityConfig{
			TLSSkipVerify: false,
			Timeout:       30 * time.Second,
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  defaultLogFile,
		},
	}
}

