package sleeper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Unit tests live in *_unit_test.go files and are named TestUnit*, so they can
// be run without the network:
//
//	go test -run '^TestUnit' ./...
//
// Every unit test client is built by newUnitTestClient, whose transport only
// talks to a local httptest server and fails any other request.

// headerOriginalHost records the host a request was addressed to before the
// unit test transport redirected it to the local test server.
const headerOriginalHost = "X-Original-Host"

// redirectedHosts are hosts the client calls directly, bypassing baseURL
// (avatars and the hidden endpoints). Requests to them are sent to the test
// server instead.
var redirectedHosts = map[string]bool{
	"api.sleeper.com": true,
	"sleepercdn.com":  true,
}

// redirectTransport sends requests for the test server and redirectedHosts to
// the test server, and refuses everything else so a unit test can never reach
// the network.
type redirectTransport struct {
	target *url.URL
	next   http.RoundTripper
}

func (rt *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch {
	case req.URL.Host == rt.target.Host:
	case redirectedHosts[req.URL.Hostname()]:
		original := req.URL.Host
		req = req.Clone(req.Context())
		req.Header.Set(headerOriginalHost, original)
		req.URL.Scheme = rt.target.Scheme
		req.URL.Host = rt.target.Host
		req.Host = ""
	default:
		return nil, fmt.Errorf("unit test attempted a network call to %s", req.URL)
	}

	return rt.next.RoundTrip(req)
}

// newUnitTestClient returns a Client built from config whose baseURL points at
// a local test server running handler.
func newUnitTestClient(t *testing.T, config Config, handler http.Handler) *Client {
	t.Helper()

	return newUnitTestClientForServer(t, config, httptest.NewServer(handler))
}

// newUnitTestClientForServer returns a Client built from config whose baseURL
// points at srv. The server is closed when the test ends.
func newUnitTestClientForServer(t *testing.T, config Config, srv *httptest.Server) *Client {
	t.Helper()

	t.Cleanup(srv.Close)

	client, err := NewClient(context.Background(), config)
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	target, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parsing test server url: %v", err)
	}

	httpClient := srv.Client()
	httpClient.Timeout = client.client.Timeout
	httpClient.Transport = &redirectTransport{target: target, next: httpClient.Transport}

	client.client = httpClient
	client.baseURL = srv.URL + "/" + string(APIVersion1)

	return client
}

// respond returns a handler that checks the request path and writes status
// and body. An empty wantPath means the client must not make a request.
func respond(t *testing.T, wantPath string, status int, body string) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		if wantPath == "" {
			t.Errorf("unexpected request to %s", r.URL)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("expected path %q, got %q", wantPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}
}

// respondExact is respond that also checks the host the client addressed and
// the raw query string. An empty wantHost means the request must have gone to
// baseURL rather than one of the redirectedHosts.
func respondExact(t *testing.T, wantHost, wantPath, wantQuery string, status int, body string) http.HandlerFunc {
	t.Helper()

	next := respond(t, wantPath, status, body)
	return func(w http.ResponseWriter, r *http.Request) {
		if wantPath != "" {
			if got := r.Header.Get(headerOriginalHost); got != wantHost {
				t.Errorf("expected request addressed to host %q, got %q", wantHost, got)
			}
			if r.URL.RawQuery != wantQuery {
				t.Errorf("expected query %q, got %q", wantQuery, r.URL.RawQuery)
			}
		}
		next(w, r)
	}
}
