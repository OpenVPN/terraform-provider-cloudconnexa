package cloudconnexa

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// retryConfig holds configuration for the retry wrapper.
type retryConfig struct {
	MaxRetries   int
	BaseDelay    time.Duration
	MaxDelay     time.Duration
	JitterFactor float64
}

// defaultRetryConfig returns sensible defaults for rate-limit retry.
func defaultRetryConfig() retryConfig {
	return retryConfig{
		MaxRetries:   5,
		BaseDelay:    1 * time.Second,
		MaxDelay:     30 * time.Second,
		JitterFactor: 0.2,
	}
}

// isRetryable returns true if the error represents a transient failure that
// should be retried (HTTP 429 Too Many Requests or 503 Service Unavailable).
func isRetryable(err error) bool {
	var apiErr *cloudconnexa.ErrClientResponse
	if errors.As(err, &apiErr) {
		code := apiErr.StatusCode()
		return code == http.StatusTooManyRequests || code == http.StatusServiceUnavailable
	}
	return false
}

// backoffDelay calculates the delay for a given attempt using exponential
// backoff with jitter, capped at cfg.MaxDelay.
func backoffDelay(cfg retryConfig, attempt int) time.Duration {
	delay := float64(cfg.BaseDelay) * math.Pow(2, float64(attempt))
	if delay > float64(cfg.MaxDelay) {
		delay = float64(cfg.MaxDelay)
	}
	// Apply jitter: delay ± jitterFactor*delay
	jitter := (rand.Float64()*2 - 1) * cfg.JitterFactor * delay
	delay += jitter
	if delay < 0 {
		delay = 0
	}
	return time.Duration(delay)
}

// withRetry executes fn and retries on retryable errors with exponential backoff.
// It respects context cancellation.
func withRetry[T any](ctx context.Context, cfg retryConfig, fn func() (T, error)) (T, error) {
	var zero T
	result, err := fn()
	if err == nil {
		return result, nil
	}
	if !isRetryable(err) {
		return zero, err
	}

	for attempt := 0; attempt < cfg.MaxRetries; attempt++ {
		delay := backoffDelay(cfg, attempt)
		tflog.Warn(ctx, "retrying after rate limit or transient error",
			map[string]interface{}{
				"attempt":   attempt + 1,
				"max":       cfg.MaxRetries,
				"delay_ms":  delay.Milliseconds(),
				"error_msg": err.Error(),
			})

		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(delay):
		}

		result, err = fn()
		if err == nil {
			return result, nil
		}
		if !isRetryable(err) {
			return zero, err
		}
	}
	return zero, err
}

// withRetryNoBody executes fn (which returns only an error) and retries on
// retryable errors with exponential backoff. Used for Update/Delete operations.
func withRetryNoBody(ctx context.Context, cfg retryConfig, fn func() error) error {
	_, err := withRetry(ctx, cfg, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	return err
}
