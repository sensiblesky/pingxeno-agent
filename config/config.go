package config

import (
	"time"
)

// Config represents the agent configuration
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Collection CollectionConfig `mapstructure:"collection"`
	Sender     SenderConfig     `mapstructure:"sender"`
	Security   SecurityConfig   `mapstructure:"security"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

// ServerConfig contains server connection settings
type ServerConfig struct {
	APIURL    string `mapstructure:"api_url"`
	APIKey    string `mapstructure:"api_key"`
	ServerKey string `mapstructure:"server_key"`
}

// CollectionConfig contains metric collection settings
type CollectionConfig struct {
	Interval time.Duration `mapstructure:"interval"`
	Jitter   time.Duration `mapstructure:"jitter"`
}

// SenderConfig contains sending/batching settings
type SenderConfig struct {
	BatchSize     int           `mapstructure:"batch_size"`
	BatchTimeout  time.Duration `mapstructure:"batch_timeout"`
	RetryAttempts int           `mapstructure:"retry_attempts"`
	RetryBackoff  time.Duration `mapstructure:"retry_backoff"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	TLSSkipVerify bool          `mapstructure:"tls_skip_verify"`
	Timeout       time.Duration `mapstructure:"timeout"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

