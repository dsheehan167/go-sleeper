package sleeper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type apiVersion string

const (
	APIVersion1 apiVersion = "v1"
)

// Client represents a client connection to the Sleeper API.
// A rate limiter is included to help prevent exceeding the API's usage limits.
// According to the Sleeper API documentation, making more than 1000 requests per minute
// may result in your IP being blocked. The rate limiter enforces a safe request rate.
type Client struct {
	client      *http.Client
	baseURL     string
	rateLimiter *rate.Limiter
}

// Config holds configuration options for creating a Client.
type Config struct {
	APIVersion     apiVersion
	Timeout        time.Duration
	RateLimitRPS   float64 // Requests per second (default: 15)
	RateLimitBurst int     // Burst capacity (default: 30)
}

const (
	defaultTimeout      = 30   // Default timeout in seconds
	defaultRateLimitRPS = 15.0 // 15 requests per second
	defaultBurst        = 30   // Allow burst of 30 requests
)

// NewClient creates a new Client using the provided configuration.
func NewClient(ctx context.Context, config Config) (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config.setDefaults()
	client := &http.Client{
		Timeout: config.Timeout * time.Second,
	}

	limit := rate.Every(time.Second / time.Duration(config.RateLimitRPS))
	rateLimiter := rate.NewLimiter(limit, config.RateLimitBurst)

	return &Client{
		client:      client,
		baseURL:     fmt.Sprintf("%s/%s", endpointBaseURL, config.APIVersion),
		rateLimiter: rateLimiter,
	}, nil
}

// validate checks if the Config fields are valid.
func (c *Config) validate() error {
	var errs []string
	if c.Timeout < 0 {
		errs = append(errs, "Timeout must be greater than zero")
	}
	if c.RateLimitRPS < 0 {
		errs = append(errs, "RateLimitRPS must be greater than or equal to zero")
	}

	if c.RateLimitBurst < 0 {
		errs = append(errs, "RateLimitBurst must be greater than or equal to zero")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

// setDefaults applies default values to zero-valued config fields
func (c *Config) setDefaults() {
	if c.Timeout == 0 {
		c.Timeout = defaultTimeout * time.Second
	}
	if c.APIVersion == "" {
		c.APIVersion = APIVersion1
	}

	if c.RateLimitRPS == 0 {
		c.RateLimitRPS = defaultRateLimitRPS
	}
	if c.RateLimitBurst == 0 {
		c.RateLimitBurst = defaultBurst
	}
}
