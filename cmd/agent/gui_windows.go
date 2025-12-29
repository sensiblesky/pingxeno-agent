//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/pingxeno/agent/agent"
	"github.com/pingxeno/agent/config"
	"github.com/pingxeno/agent/internal/logger"
	"go.uber.org/zap"
)

// runGUI starts the Windows GUI application
func runGUI() error {
	// Check if we're on Windows
	if runtime.GOOS != "windows" {
		return fmt.Errorf("GUI is only available on Windows")
	}

	// Load config
	configPath := findConfigFile()
	if configPath == "" {
		// Try default Windows location
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = "C:\\ProgramData"
		}
		configPath = filepath.Join(programData, "PingXeno", "agent.yaml")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		cfg = config.DefaultConfig()
		// Try to save default config
		os.MkdirAll(filepath.Dir(configPath), 0755)
		config.SaveConfig(cfg, configPath)
	}

	// Setup logger with file output
	log, err := logger.NewLogger(cfg.Logging.Level, cfg.Logging.File)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer log.Sync()

	// Create agent
	agentInstance, err := agent.NewAgent(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}

	// Run agent in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Received shutdown signal")
		cancel()
	}()

	log.Info("Agent started in background mode",
		zap.String("log_file", cfg.Logging.File),
		zap.String("config_file", configPath),
		zap.String("api_url", cfg.Server.APIURL),
	)

	// Run the agent (this will block until context is cancelled)
	// All output goes to log file, so it can run in background
	return agentInstance.Run(ctx)
}

