package sleeper

import (
	"context"
	"testing"
)

func TestClientNew(t *testing.T) {
	tt := []struct {
		testcase   string
		config     Config
		shouldPass bool
	}{
		{
			"default config",
			Config{},
			true,
		},
		{
			"custom config",
			Config{
				APIVersion:   APIVersion1,
				Timeout:      10,
				RateLimitRPS: 5,
			},
			true,
		},
		{
			"invalid rate limit",
			Config{
				Timeout:        -1,
				RateLimitRPS:   -1,
				RateLimitBurst: -1,
			},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client, err := NewClient(context.Background(), tc.config)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got client: %v", client)
				return
			}

			if client == nil {
				t.Errorf("expected non-nil client")
				return
			}
		})
	}
}
