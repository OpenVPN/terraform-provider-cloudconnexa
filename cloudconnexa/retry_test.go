package cloudconnexa

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"github.com/stretchr/testify/require"
)

// newRetryTestClient creates a client pointing at a test server.
func newRetryTestClient(t *testing.T, handler http.Handler) *cloudconnexa.Client {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"unit-test-token"}`))
	})
	mux.Handle("/", handler)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	c, err := cloudconnexa.NewClientWithOptions(server.URL, "test-id", "test-secret", &cloudconnexa.ClientOptions{
		AllowInsecureHTTP: true,
	})
	require.NoError(t, err)
	return c
}

// makeAPIError triggers a real SDK error with the given status code.
func makeAPIError(t *testing.T, statusCode int) error {
	t.Helper()
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`{"error":"test"}`))
	})
	c := newRetryTestClient(t, handler)
	_, err := c.Hosts.Get("nonexistent")
	require.Error(t, err)
	return err
}

func TestUnit_isRetryable_429(t *testing.T) {
	err := makeAPIError(t, http.StatusTooManyRequests)
	if !isRetryable(err) {
		t.Error("expected 429 to be retryable")
	}
}

func TestUnit_isRetryable_503(t *testing.T) {
	err := makeAPIError(t, http.StatusServiceUnavailable)
	if !isRetryable(err) {
		t.Error("expected 503 to be retryable")
	}
}

func TestUnit_isRetryable_400(t *testing.T) {
	err := makeAPIError(t, http.StatusBadRequest)
	if isRetryable(err) {
		t.Error("expected 400 to NOT be retryable")
	}
}

func TestUnit_isRetryable_404(t *testing.T) {
	err := makeAPIError(t, http.StatusNotFound)
	if isRetryable(err) {
		t.Error("expected 404 to NOT be retryable")
	}
}

func TestUnit_isRetryable_GenericError(t *testing.T) {
	err := errors.New("connection reset")
	if isRetryable(err) {
		t.Error("expected generic error to NOT be retryable")
	}
}

func TestUnit_backoffDelay_ExponentialGrowth(t *testing.T) {
	cfg := retryConfig{
		MaxRetries:   5,
		BaseDelay:    1 * time.Second,
		MaxDelay:     30 * time.Second,
		JitterFactor: 0,
	}

	expected := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		16 * time.Second,
	}

	for attempt, want := range expected {
		got := backoffDelay(cfg, attempt)
		if got != want {
			t.Errorf("attempt %d: got %v, want %v", attempt, got, want)
		}
	}
}

func TestUnit_backoffDelay_CappedAtMaxDelay(t *testing.T) {
	cfg := retryConfig{
		MaxRetries:   10,
		BaseDelay:    1 * time.Second,
		MaxDelay:     5 * time.Second,
		JitterFactor: 0,
	}

	got := backoffDelay(cfg, 5)
	if got != 5*time.Second {
		t.Errorf("expected delay capped at 5s, got %v", got)
	}
}

func TestUnit_backoffDelay_WithJitter(t *testing.T) {
	cfg := retryConfig{
		MaxRetries:   5,
		BaseDelay:    1 * time.Second,
		MaxDelay:     30 * time.Second,
		JitterFactor: 0.2,
	}

	for i := 0; i < 100; i++ {
		got := backoffDelay(cfg, 0)
		if got < 800*time.Millisecond || got > 1200*time.Millisecond {
			t.Errorf("attempt 0 with jitter: got %v, expected [800ms, 1200ms]", got)
		}
	}
}

func TestUnit_withRetry_SuccessOnFirstCall(t *testing.T) {
	cfg := defaultRetryConfig()
	ctx := context.Background()

	calls := 0
	result, err := withRetry(ctx, cfg, func() (string, error) {
		calls++
		return "ok", nil
	})

	require.NoError(t, err)
	require.Equal(t, "ok", result)
	require.Equal(t, 1, calls)
}

func TestUnit_withRetry_RetriesOn429ThenSucceeds(t *testing.T) {
	var callCount int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"rate limited"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"host-1","name":"test","description":"","connectors":[]}`))
	})
	c := newRetryTestClient(t, handler)

	cfg := retryConfig{
		MaxRetries:   5,
		BaseDelay:    1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		JitterFactor: 0,
	}
	ctx := context.Background()

	host, err := withRetry(ctx, cfg, func() (*cloudconnexa.Host, error) {
		return c.Hosts.Get("host-1")
	})

	require.NoError(t, err)
	require.NotNil(t, host)
	require.Equal(t, "host-1", host.ID)
	require.Equal(t, int32(3), atomic.LoadInt32(&callCount))
}

func TestUnit_withRetry_ExhaustsRetries(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"rate limited"}`))
	})
	c := newRetryTestClient(t, handler)

	cfg := retryConfig{
		MaxRetries:   2,
		BaseDelay:    1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		JitterFactor: 0,
	}
	ctx := context.Background()

	var callCount int32
	_, err := withRetry(ctx, cfg, func() (*cloudconnexa.Host, error) {
		atomic.AddInt32(&callCount, 1)
		return c.Hosts.Get("host-1")
	})

	require.Error(t, err)
	require.Equal(t, int32(3), atomic.LoadInt32(&callCount))
}

func TestUnit_withRetry_NonRetryableErrorReturnsImmediately(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	})
	c := newRetryTestClient(t, handler)

	cfg := retryConfig{
		MaxRetries:   5,
		BaseDelay:    1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		JitterFactor: 0,
	}
	ctx := context.Background()

	var callCount int32
	_, err := withRetry(ctx, cfg, func() (*cloudconnexa.Host, error) {
		atomic.AddInt32(&callCount, 1)
		return c.Hosts.Get("host-1")
	})

	require.Error(t, err)
	require.Equal(t, int32(1), atomic.LoadInt32(&callCount))
}

func TestUnit_withRetry_ContextCancellation(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"rate limited"}`))
	})
	c := newRetryTestClient(t, handler)

	cfg := retryConfig{
		MaxRetries:   10,
		BaseDelay:    100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		JitterFactor: 0,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := withRetry(ctx, cfg, func() (*cloudconnexa.Host, error) {
		return c.Hosts.Get("host-1")
	})

	require.ErrorIs(t, err, context.Canceled)
}

func TestUnit_withRetryNoBody_Success(t *testing.T) {
	var callCount int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"rate limited"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	c := newRetryTestClient(t, handler)

	cfg := retryConfig{
		MaxRetries:   3,
		BaseDelay:    1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		JitterFactor: 0,
	}
	ctx := context.Background()

	err := withRetryNoBody(ctx, cfg, func() error {
		return c.Hosts.Delete("host-1")
	})

	require.NoError(t, err)
	require.Equal(t, int32(2), atomic.LoadInt32(&callCount))
}

func TestUnit_withRetry_503ThenSucceeds(t *testing.T) {
	var callCount int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"service unavailable"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"host-1","name":"test","description":"","connectors":[]}`))
	})
	c := newRetryTestClient(t, handler)

	cfg := retryConfig{
		MaxRetries:   3,
		BaseDelay:    1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		JitterFactor: 0,
	}
	ctx := context.Background()

	host, err := withRetry(ctx, cfg, func() (*cloudconnexa.Host, error) {
		return c.Hosts.Get("host-1")
	})

	require.NoError(t, err)
	require.NotNil(t, host)
	require.Equal(t, "host-1", host.ID)
	require.Equal(t, int32(2), atomic.LoadInt32(&callCount))
}
