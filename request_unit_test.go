package sleeper

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// blockUntilDone holds the response until the client gives up on the request.
func blockUntilDone(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
	case <-time.After(5 * time.Second):
		w.WriteHeader(http.StatusOK)
	}
}

func TestUnitClient_GetRequest(t *testing.T) {
	tt := []struct {
		testcase   string
		handler    http.HandlerFunc
		expected   string
		wantStatus int
		wantErr    string
		shouldPass bool
	}{
		{
			"200 returns body",
			func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, `{"ok":true}`)
			},
			`{"ok":true}`,
			0,
			"",
			true,
		},
		{
			"201 is treated as success",
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = io.WriteString(w, `created`)
			},
			"created",
			0,
			"",
			true,
		},
		{
			"204 returns empty body",
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			"",
			0,
			"",
			true,
		},
		{
			"304 is an api error",
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotModified)
			},
			"",
			http.StatusNotModified,
			"API error: 304 Not Modified",
			false,
		},
		{
			"400 is an api error",
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			"",
			http.StatusBadRequest,
			"API error: 400 Bad Request",
			false,
		},
		{
			"404 is an api error",
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			"",
			http.StatusNotFound,
			"API error: 404 Not Found",
			false,
		},
		{
			"429 is an api error",
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
			},
			"",
			http.StatusTooManyRequests,
			"API error: 429 Too Many Requests",
			false,
		},
		{
			"500 is an api error",
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, `{"error":"boom"}`)
			},
			"",
			http.StatusInternalServerError,
			"API error: 500 Internal Server Error",
			false,
		},
		{
			"truncated body is a read error",
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Length", "100")
				_, _ = io.WriteString(w, `{"partial"`)
			},
			"",
			0,
			"reading response",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, tc.handler)

			by, err := client.getRequest(t.Context(), client.baseURL+"/test")
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("expected error containing %q, got %q", tc.wantErr, err)
				}
				var apiErr *APIError
				if tc.wantStatus != 0 {
					if !errors.As(err, &apiErr) {
						t.Errorf("expected *APIError, got %T", err)
						return
					}
					if apiErr.StatusCode != tc.wantStatus {
						t.Errorf("expected status %d, got %d", tc.wantStatus, apiErr.StatusCode)
					}
				} else if errors.As(err, &apiErr) {
					t.Errorf("expected a non-API error, got %v", err)
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got body: %q", by)
				return
			}

			if string(by) != tc.expected {
				t.Errorf("expected body %q, got %q", tc.expected, by)
				return
			}
		})
	}
}

func TestUnitClient_GetRequestSetup(t *testing.T) {
	tt := []struct {
		testcase   string
		setup      func(t *testing.T) (*Client, string)
		wantErr    string
		shouldPass bool
	}{
		{
			"nil http client",
			func(_ *testing.T) (*Client, string) {
				return &Client{}, "/test"
			},
			"http client is not initialised",
			false,
		},
		{
			"invalid endpoint url",
			func(t *testing.T) (*Client, string) {
				return newUnitTestClient(t, Config{}, respond(t, "", 0, "")), "http://[::1"
			},
			"creating request",
			false,
		},
		{
			"server unreachable",
			func(t *testing.T) (*Client, string) {
				srv := httptest.NewServer(respond(t, "", 0, ""))
				client := newUnitTestClientForServer(t, Config{}, srv)
				srv.Close()
				return client, client.baseURL + "/test"
			},
			"completing request",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client, endpoint := tc.setup(t)

			by, err := client.getRequest(t.Context(), endpoint)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("expected error containing %q, got %q", tc.wantErr, err)
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got body: %q", by)
				return
			}
		})
	}
}

func TestUnitClient_GetRequestCancellation(t *testing.T) {
	tt := []struct {
		testcase      string
		clientTimeout time.Duration
		ctx           func() (context.Context, context.CancelFunc)
		handler       http.HandlerFunc
		wantErr       error
		wantRequests  int32
		shouldPass    bool
	}{
		{
			"context canceled before request",
			0,
			func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx, cancel
			},
			blockUntilDone,
			context.Canceled,
			0,
			false,
		},
		{
			"context canceled during request",
			0,
			func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(t.Context())
				time.AfterFunc(50*time.Millisecond, cancel)
				return ctx, cancel
			},
			blockUntilDone,
			context.Canceled,
			1,
			false,
		},
		{
			"context deadline exceeded during request",
			0,
			func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(t.Context(), 50*time.Millisecond)
			},
			blockUntilDone,
			context.DeadlineExceeded,
			1,
			false,
		},
		{
			"http client timeout",
			50 * time.Millisecond,
			func() (context.Context, context.CancelFunc) {
				return context.WithCancel(t.Context())
			},
			blockUntilDone,
			context.DeadlineExceeded,
			1,
			false,
		},
		{
			"response within deadline",
			time.Second,
			func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(t.Context(), time.Second)
			},
			func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, `{}`)
			},
			nil,
			1,
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			var requests atomic.Int32
			client := newUnitTestClient(t, Config{Timeout: tc.clientTimeout}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				tc.handler(w, r)
			}))

			ctx, cancel := tc.ctx()
			defer cancel()

			start := time.Now()
			_, err := client.getRequest(ctx, client.baseURL+"/slow")
			elapsed := time.Since(start)

			if got := requests.Load(); got != tc.wantRequests {
				t.Errorf("expected %d requests to reach the server, got %d", tc.wantRequests, got)
			}

			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error wrapping %v, got %v", tc.wantErr, err)
				}
				if elapsed > 2*time.Second {
					t.Errorf("expected request to be abandoned promptly, took %v", elapsed)
				}
				t.Logf("expected error after %v: %v", elapsed, err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but request succeeded")
				return
			}
		})
	}
}

func TestUnitClient_RateLimit(t *testing.T) {
	tt := []struct {
		testcase       string
		config         Config
		requests       int
		lastTimeout    time.Duration
		wantMinElapsed time.Duration
		wantServed     int32
		shouldPass     bool
	}{
		{
			"burst serves requests without waiting",
			Config{RateLimitRPS: 1, RateLimitBurst: 3},
			3,
			0,
			0,
			3,
			true,
		},
		{
			"requests beyond burst wait for tokens",
			Config{RateLimitRPS: 20, RateLimitBurst: 1},
			3,
			0,
			90 * time.Millisecond, // two waits of 50ms each, with scheduling slack
			3,
			true,
		},
		{
			"wait longer than context deadline fails without a request",
			Config{RateLimitRPS: 1, RateLimitBurst: 1},
			2,
			100 * time.Millisecond,
			0,
			1,
			false,
		},
		{
			"exhausted burst fails without a request",
			Config{RateLimitRPS: 1, RateLimitBurst: 3},
			4,
			100 * time.Millisecond,
			0,
			3,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			var served atomic.Int32
			client := newUnitTestClient(t, tc.config, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				served.Add(1)
				_, _ = io.WriteString(w, `{}`)
			}))

			start := time.Now()
			var err error
			for i := range tc.requests {
				ctx, cancel := t.Context(), context.CancelFunc(func() {})
				if i == tc.requests-1 && tc.lastTimeout > 0 {
					ctx, cancel = context.WithTimeout(t.Context(), tc.lastTimeout)
				}
				_, err = client.getRequest(ctx, client.baseURL+"/limited")
				cancel()
				if err != nil {
					break
				}
			}
			elapsed := time.Since(start)

			if got := served.Load(); got != tc.wantServed {
				t.Errorf("expected %d requests served, got %d", tc.wantServed, got)
			}

			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected rate limiter to reject the last request")
				return
			}

			if elapsed < tc.wantMinElapsed {
				t.Errorf("expected requests to take at least %v, took %v", tc.wantMinElapsed, elapsed)
				return
			}
		})
	}
}
