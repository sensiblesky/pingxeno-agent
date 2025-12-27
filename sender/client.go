package sender

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pingxeno/agent/config"
	"github.com/pingxeno/agent/protocol"
	"go.uber.org/zap"
)

// Client handles HTTP communication with the API
type Client struct {
	config     *config.Config
	httpClient *http.Client
	logger     *zap.Logger
}

// NewClient creates a new API client
func NewClient(cfg *config.Config, logger *zap.Logger) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.Security.TLSSkipVerify,
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   cfg.Security.Timeout,
	}

	return &Client{
		config:     cfg,
		httpClient:  client,
		logger:     logger,
	}
}

// SendMetrics sends metrics payload to the API
func (c *Client) SendMetrics(payload *protocol.MetricsPayload) error {
	payload.ServerKey = c.config.Server.ServerKey

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.Server.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.config.Server.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned error: %d - %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err == nil {
		c.logger.Info("Metrics sent successfully",
			zap.String("server_id", fmt.Sprintf("%v", response["server_id"])),
			zap.String("stat_id", fmt.Sprintf("%v", response["stat_id"])),
		)
	}

	return nil
}

// TestConnection tests the connection to the API
func (c *Client) TestConnection() error {
	// Create a minimal test payload
	testPayload := &protocol.MetricsPayload{
		ServerKey:  c.config.Server.ServerKey,
		RecordedAt: time.Now(),
	}

	return c.SendMetrics(testPayload)
}

