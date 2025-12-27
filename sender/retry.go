package sender

import (
	"fmt"
	"time"

	"github.com/pingxeno/agent/protocol"
	"go.uber.org/zap"
)

// RetrySender wraps a client with retry logic
type RetrySender struct {
	client *Client
	config *struct {
		Attempts int
		Backoff  time.Duration
	}
	logger *zap.Logger
}

// NewRetrySender creates a new retry sender
func NewRetrySender(client *Client, attempts int, backoff time.Duration, logger *zap.Logger) *RetrySender {
	return &RetrySender{
		client: client,
		config: &struct {
			Attempts int
			Backoff  time.Duration
		}{
			Attempts: attempts,
			Backoff:  backoff,
		},
		logger: logger,
	}
}

// SendWithRetry sends metrics with retry logic
func (r *RetrySender) SendWithRetry(payload *protocol.MetricsPayload) error {
	var lastErr error

	for attempt := 1; attempt <= r.config.Attempts; attempt++ {
		err := r.client.SendMetrics(payload)
		if err == nil {
			if attempt > 1 {
				r.logger.Info("Successfully sent after retry",
					zap.Int("attempt", attempt),
				)
			}
			return nil
		}

		lastErr = err
		r.logger.Warn("Failed to send metrics, retrying",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", r.config.Attempts),
			zap.Error(err),
		)

		if attempt < r.config.Attempts {
			backoff := time.Duration(attempt) * r.config.Backoff
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("failed to send after %d attempts: %w", r.config.Attempts, lastErr)
}

