package config

import "time"

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
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
			File:  "",
		},
	}
}

