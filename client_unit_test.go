package sleeper

import (
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestUnitConfig_Validate(t *testing.T) {
	tt := []struct {
		testcase   string
		config     Config
		wantErrs   []string
		shouldPass bool
	}{
		{
			"zero config",
			Config{},
			nil,
			true,
		},
		{
			"positive values",
			Config{Timeout: time.Second, RateLimitRPS: 5, RateLimitBurst: 10},
			nil,
			true,
		},
		{
			"negative timeout",
			Config{Timeout: -time.Second},
			[]string{"Timeout"},
			false,
		},
		{
			"negative rate limit rps",
			Config{RateLimitRPS: -1},
			[]string{"RateLimitRPS"},
			false,
		},
		{
			"negative rate limit burst",
			Config{RateLimitBurst: -1},
			[]string{"RateLimitBurst"},
			false,
		},
		{
			"all fields invalid",
			Config{Timeout: -1, RateLimitRPS: -1, RateLimitBurst: -1},
			[]string{"Timeout", "RateLimitRPS", "RateLimitBurst"},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			err := tc.config.validate()
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				for _, want := range tc.wantErrs {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("expected error to mention %q, got %q", want, err)
					}
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure for config: %+v", tc.config)
				return
			}
		})
	}
}

func TestUnitConfig_SetDefaults(t *testing.T) {
	tt := []struct {
		testcase string
		config   Config
		expected Config
	}{
		{
			"zero config gets defaults",
			Config{},
			Config{
				APIVersion:     APIVersion1,
				Timeout:        defaultTimeout * time.Second,
				RateLimitRPS:   defaultRateLimitRPS,
				RateLimitBurst: defaultBurst,
			},
		},
		{
			"custom values are kept",
			Config{
				APIVersion:     APIVersion1,
				Timeout:        5 * time.Second,
				RateLimitRPS:   2,
				RateLimitBurst: 4,
			},
			Config{
				APIVersion:     APIVersion1,
				Timeout:        5 * time.Second,
				RateLimitRPS:   2,
				RateLimitBurst: 4,
			},
		},
		{
			"only zero fields are defaulted",
			Config{Timeout: 5 * time.Second, RateLimitBurst: 4},
			Config{
				APIVersion:     APIVersion1,
				Timeout:        5 * time.Second,
				RateLimitRPS:   defaultRateLimitRPS,
				RateLimitBurst: 4,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			config := tc.config
			config.setDefaults()

			if config != tc.expected {
				t.Errorf("expected config %+v, got %+v", tc.expected, config)
				return
			}
		})
	}
}

func TestUnitClient_New(t *testing.T) {
	tt := []struct {
		testcase      string
		config        Config
		expectedURL   string
		expectedTime  time.Duration
		expectedLimit rate.Limit
		expectedBurst int
		shouldPass    bool
	}{
		{
			"default config",
			Config{},
			endpointBaseURL + "/v1",
			defaultTimeout * time.Second,
			rate.Every(time.Second / 15), // NewClient truncates the interval to whole nanoseconds
			defaultBurst,
			true,
		},
		{
			"custom config",
			Config{
				APIVersion:     APIVersion1,
				Timeout:        2 * time.Second,
				RateLimitRPS:   4,
				RateLimitBurst: 8,
			},
			endpointBaseURL + "/v1",
			2 * time.Second,
			rate.Every(250 * time.Millisecond),
			8,
			true,
		},
		{
			"invalid config",
			Config{Timeout: -1},
			"",
			0,
			0,
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client, err := NewClient(t.Context(), tc.config)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !strings.HasPrefix(err.Error(), "invalid config") {
					t.Errorf("expected error to start with %q, got %q", "invalid config", err)
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got client: %v", client)
				return
			}

			if client.baseURL != tc.expectedURL {
				t.Errorf("expected baseURL %q, got %q", tc.expectedURL, client.baseURL)
			}
			if client.client.Timeout != tc.expectedTime {
				t.Errorf("expected timeout %v, got %v", tc.expectedTime, client.client.Timeout)
			}
			if client.rateLimiter.Limit() != tc.expectedLimit {
				t.Errorf("expected rate limit %v, got %v", tc.expectedLimit, client.rateLimiter.Limit())
			}
			if client.rateLimiter.Burst() != tc.expectedBurst {
				t.Errorf("expected burst %d, got %d", tc.expectedBurst, client.rateLimiter.Burst())
			}
		})
	}
}
