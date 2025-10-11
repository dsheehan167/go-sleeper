package sleeper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error response from the Sleeper API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d %s", e.StatusCode, e.Message)
}

func (c *Client) getRequest(ctx context.Context, endpoint string) ([]byte, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("http client cannot be nil")
	}

	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("completing request: %w", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Return successful responses
	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
		return b, nil
	}

	return nil, &APIError{
		StatusCode: res.StatusCode,
		Message:    http.StatusText(res.StatusCode),
	}
}
