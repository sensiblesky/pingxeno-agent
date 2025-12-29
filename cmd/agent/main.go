package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pingxeno/agent/agent"
	"github.com/pingxeno/agent/config"
	"github.com/pingxeno/agent/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	version = "1.0.0"
	cfg     *config.Config
	log     *zap.Logger
)

func main() {
	log, _ = logger.NewLogger("info", "")

	// Check if running with --gui flag or on Windows without console
	if len(os.Args) > 1 && os.Args[1] == "gui" {
		if err := runGUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	rootCmd := &cobra.Command{
		Use:   "pingxeno-agent",
		Short: "PingXeno Monitoring Agent",
		Long:  "A cross-platform system monitoring agent for PingXeno",
	}

	rootCmd.AddCommand(
		createRunCommand(),
		createInstallCommand(),
		createUninstallCommand(),
		createStatusCommand(),
		createConfigCommand(),
		createTestCommand(),
		createVersionCommand(),
		createGUICommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func createRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the agent",
		Long:  "Start the monitoring agent and begin collecting metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, _ := cmd.Flags().GetString("config")
			if configPath == "" {
				configPath = findConfigFile()
			}

			var err error
			cfg, err = config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			logFile, _ := cmd.Flags().GetString("log-file")
			if logFile != "" {
				cfg.Logging.File = logFile
			} else if cfg.Logging.File == "" {
				// Use default log file if not specified
				cfg.Logging.File = getDefaultLogFile()
			}

			log, err = logger.NewLogger(cfg.Logging.Level, cfg.Logging.File)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}
			defer log.Sync()

			agent, err := agent.NewAgent(cfg, log)
			if err != nil {
				return fmt.Errorf("failed to create agent: %w", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle signals
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-sigChan
				log.Info("Received shutdown signal")
				cancel()
			}()

			return agent.Run(ctx)
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to configuration file")
	cmd.Flags().String("log-file", "", "Path to log file")

	return cmd
}

func createInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install and configure the agent",
		Long:  "Interactive installation that prompts for configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("=== PingXeno Agent Installation ===")
			fmt.Println()

			reader := bufio.NewReader(os.Stdin)

			// Get API URL
			fmt.Print("Enter API URL (e.g., https://your-domain.com/api/v1/server-stats): ")
			apiURL, _ := reader.ReadString('\n')
			apiURL = strings.TrimSpace(apiURL)
			if apiURL == "" {
				return fmt.Errorf("API URL is required")
			}

			// Get API Key
			fmt.Print("Enter API Key: ")
			apiKey, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)
			if apiKey == "" {
				return fmt.Errorf("API Key is required")
			}

			// Get Server Key
			fmt.Print("Enter Server Key: ")
			serverKey, _ := reader.ReadString('\n')
			serverKey = strings.TrimSpace(serverKey)
			if serverKey == "" {
				return fmt.Errorf("Server Key is required")
			}

			// Get Collection Interval
			fmt.Print("Enter collection interval in seconds (default: 60): ")
			intervalStr, _ := reader.ReadString('\n')
			intervalStr = strings.TrimSpace(intervalStr)
			interval := 60 * time.Second
			if intervalStr != "" {
				var seconds int
				if _, err := fmt.Sscanf(intervalStr, "%d", &seconds); err == nil && seconds > 0 {
					interval = time.Duration(seconds) * time.Second
				}
			}

			// Create config
			cfg = config.DefaultConfig()
			cfg.Server.APIURL = apiURL
			cfg.Server.APIKey = apiKey
			cfg.Server.ServerKey = serverKey
			cfg.Collection.Interval = interval

			// Determine config path
			configPath, _ := cmd.Flags().GetString("config")
			if configPath == "" {
				configPath = getDefaultConfigPath()
			}

			// Save config
			if err := config.SaveConfig(cfg, configPath); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("\n✓ Configuration saved to: %s\n", configPath)
			fmt.Println("\nYou can now run the agent with:")
			fmt.Printf("  pingxeno-agent run --config %s\n", configPath)

			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to save configuration file")

	return cmd
}

func createConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Reconfigure the agent",
		Long:  "Interactive configuration update",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, _ := cmd.Flags().GetString("config")
			if configPath == "" {
				configPath = findConfigFile()
			}

			var err error
			cfg, err = config.LoadConfig(configPath)
			if err != nil {
				cfg = config.DefaultConfig()
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Println("=== PingXeno Agent Configuration ===")
			fmt.Println("Press Enter to keep current value")
			fmt.Println()

			// API URL
			fmt.Printf("API URL [%s]: ", cfg.Server.APIURL)
			apiURL, _ := reader.ReadString('\n')
			apiURL = strings.TrimSpace(apiURL)
			if apiURL != "" {
				cfg.Server.APIURL = apiURL
			}

			// API Key
			fmt.Printf("API Key [***hidden***]: ")
			apiKey, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)
			if apiKey != "" {
				cfg.Server.APIKey = apiKey
			}

			// Server Key
			fmt.Printf("Server Key [***hidden***]: ")
			serverKey, _ := reader.ReadString('\n')
			serverKey = strings.TrimSpace(serverKey)
			if serverKey != "" {
				cfg.Server.ServerKey = serverKey
			}

			// Collection Interval
			fmt.Printf("Collection interval in seconds [%d]: ", int(cfg.Collection.Interval.Seconds()))
			intervalStr, _ := reader.ReadString('\n')
			intervalStr = strings.TrimSpace(intervalStr)
			if intervalStr != "" {
				var seconds int
				if _, err := fmt.Sscanf(intervalStr, "%d", &seconds); err == nil && seconds > 0 {
					cfg.Collection.Interval = time.Duration(seconds) * time.Second
				}
			}

			// Save config
			if err := config.SaveConfig(cfg, configPath); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("\n✓ Configuration updated: %s\n", configPath)
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to configuration file")

	return cmd
}

func createStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check agent status",
		Long:  "Test connection and verify agent is working correctly",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, _ := cmd.Flags().GetString("config")
			if configPath == "" {
				configPath = findConfigFile()
			}

			var err error
			cfg, err = config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			log, err = logger.NewLogger("info", "")
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}
			defer log.Sync()

			fmt.Println("=== PingXeno Agent Status ===")
			fmt.Println()

			// Test connection
			agent, err := agent.NewAgent(cfg, log)
			if err != nil {
				return fmt.Errorf("failed to create agent: %w", err)
			}

			fmt.Println("Testing API connection...")
			if err := agent.TestConnection(); err != nil {
				fmt.Printf("✗ Connection failed: %v\n", err)
				return err
			}
			fmt.Println("✓ API connection successful")
			fmt.Println()

			// Collect sample metrics
			fmt.Println("Collecting system metrics...")
			payload, err := agent.CollectMetrics()
			if err != nil {
				fmt.Printf("✗ Failed to collect metrics: %v\n", err)
				return err
			}

			// Try to send metrics
			fmt.Println("Sending test metrics...")
			if err := agent.Sender().SendWithRetry(payload); err != nil {
				fmt.Printf("⚠ Warning: Failed to send test metrics: %v\n", err)
			} else {
				fmt.Println("✓ Test metrics sent successfully")
			}

			fmt.Println("✓ Metrics collected successfully")
			fmt.Println()
			fmt.Println("Sample Metrics:")
			if payload.CPUUsagePercent != nil {
				fmt.Printf("  CPU Usage: %.2f%%\n", *payload.CPUUsagePercent)
			}
			if payload.MemoryUsagePercent != nil {
				fmt.Printf("  Memory Usage: %.2f%%\n", *payload.MemoryUsagePercent)
			}
			if payload.DiskUsagePercent != nil {
				fmt.Printf("  Disk Usage: %.2f%%\n", *payload.DiskUsagePercent)
			}
			if payload.UptimeSeconds != nil {
				fmt.Printf("  Uptime: %d seconds\n", *payload.UptimeSeconds)
			}

			fmt.Println()
			fmt.Println("✓ Agent is working correctly")

			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to configuration file")

	return cmd
}

func createTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test sending metrics",
		Long:  "Collect and send a single metrics payload for testing",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, _ := cmd.Flags().GetString("config")
			if configPath == "" {
				configPath = findConfigFile()
			}

			var err error
			cfg, err = config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			log, err = logger.NewLogger("info", "")
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}
			defer log.Sync()

			agent, err := agent.NewAgent(cfg, log)
			if err != nil {
				return fmt.Errorf("failed to create agent: %w", err)
			}

			fmt.Println("Collecting metrics...")
			payload, err := agent.CollectMetrics()
			if err != nil {
				return fmt.Errorf("failed to collect metrics: %w", err)
			}

			fmt.Println("Sending metrics to API...")
			if err := agent.Sender().SendWithRetry(payload); err != nil {
				return fmt.Errorf("failed to send metrics: %w", err)
			}

			fmt.Println("✓ Metrics sent successfully!")
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to configuration file")

	return cmd
}

func createUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall the agent service",
		Long:  "Remove the agent service from the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Uninstall functionality will be implemented per platform")
			return nil
		},
	}
}

func createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("PingXeno Agent v%s\n", version)
		},
	}
}

func createGUICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gui",
		Short: "Run agent with GUI (Windows only)",
		Long:  "Start the agent in background mode with GUI interface (Windows only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGUI()
		},
	}
}

func findConfigFile() string {
	paths := []string{
		"./agent.yaml",
		"/etc/pingxeno/agent.yaml",
		os.Getenv("HOME") + "/.config/pingxeno/agent.yaml",
		os.Getenv("HOME") + "/.pingxeno/agent.yaml",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func getDefaultConfigPath() string {
	if os.Geteuid() == 0 {
		return "/etc/pingxeno/agent.yaml"
	}
	return os.Getenv("HOME") + "/.config/pingxeno/agent.yaml"
}

func getDefaultLogFile() string {
	if runtime.GOOS == "windows" {
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = "C:\\ProgramData"
		}
		return filepath.Join(programData, "PingXeno", "agent.log")
	}
	return "/var/log/pingxeno-agent.log"
}

