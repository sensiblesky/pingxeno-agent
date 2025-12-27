package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig loads configuration from file, environment variables, or flags
func LoadConfig(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Default config locations
		viper.SetConfigName("agent")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/pingxeno")
		viper.AddConfigPath("$HOME/.config/pingxeno")
		viper.AddConfigPath("$HOME/.pingxeno")
	}

	// Environment variables
	viper.SetEnvPrefix("PINGXENO")
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("server.api_url", "PINGXENO_API_URL")
	viper.BindEnv("server.api_key", "PINGXENO_API_KEY")
	viper.BindEnv("server.server_key", "PINGXENO_SERVER_KEY")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate required fields
	if cfg.Server.APIKey == "" {
		return nil, fmt.Errorf("server.api_key is required")
	}
	if cfg.Server.ServerKey == "" {
		return nil, fmt.Errorf("server.server_key is required")
	}
	if cfg.Server.APIURL == "" {
		return nil, fmt.Errorf("server.api_url is required")
	}

	// Parse duration strings from viper if they're strings
	if intervalStr := viper.GetString("collection.interval"); intervalStr != "" {
		if d, err := time.ParseDuration(intervalStr); err == nil {
			cfg.Collection.Interval = d
		}
	}
	if cfg.Collection.Interval == 0 {
		cfg.Collection.Interval = 60 * time.Second
	}

	if jitterStr := viper.GetString("collection.jitter"); jitterStr != "" {
		if d, err := time.ParseDuration(jitterStr); err == nil {
			cfg.Collection.Jitter = d
		}
	}
	if cfg.Collection.Jitter == 0 {
		cfg.Collection.Jitter = 5 * time.Second
	}

	if timeoutStr := viper.GetString("sender.batch_timeout"); timeoutStr != "" {
		if d, err := time.ParseDuration(timeoutStr); err == nil {
			cfg.Sender.BatchTimeout = d
		}
	}
	if cfg.Sender.BatchTimeout == 0 {
		cfg.Sender.BatchTimeout = 30 * time.Second
	}

	if backoffStr := viper.GetString("sender.retry_backoff"); backoffStr != "" {
		if d, err := time.ParseDuration(backoffStr); err == nil {
			cfg.Sender.RetryBackoff = d
		}
	}
	if cfg.Sender.RetryBackoff == 0 {
		cfg.Sender.RetryBackoff = 2 * time.Second
	}

	if timeoutStr := viper.GetString("security.timeout"); timeoutStr != "" {
		if d, err := time.ParseDuration(timeoutStr); err == nil {
			cfg.Security.Timeout = d
		}
	}
	if cfg.Security.Timeout == 0 {
		cfg.Security.Timeout = 30 * time.Second
	}

	return cfg, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(cfg *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Use a fresh viper instance
	v := viper.New()
	v.SetConfigType("yaml")

	v.Set("server.api_url", cfg.Server.APIURL)
	v.Set("server.api_key", cfg.Server.APIKey)
	v.Set("server.server_key", cfg.Server.ServerKey)
	v.Set("collection.interval", cfg.Collection.Interval.String())
	v.Set("collection.jitter", cfg.Collection.Jitter.String())
	v.Set("sender.batch_size", cfg.Sender.BatchSize)
	v.Set("sender.batch_timeout", cfg.Sender.BatchTimeout.String())
	v.Set("sender.retry_attempts", cfg.Sender.RetryAttempts)
	v.Set("sender.retry_backoff", cfg.Sender.RetryBackoff.String())
	v.Set("security.tls_skip_verify", cfg.Security.TLSSkipVerify)
	v.Set("security.timeout", cfg.Security.Timeout.String())
	v.Set("logging.level", cfg.Logging.Level)
	v.Set("logging.file", cfg.Logging.File)

	return v.WriteConfigAs(path)
}

