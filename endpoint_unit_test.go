package sleeper

import "testing"

func TestUnitEndpoint_Build(t *testing.T) {
	tt := []struct {
		testcase string
		template string
		args     []any
		expected string
	}{
		{
			"no args",
			endpointNFLPlayers,
			nil,
			"https://api.sleeper.app/v1/players/nfl",
		},
		{
			"string arg",
			endpointUser,
			[]any{"sleeperuser"},
			"https://api.sleeper.app/v1/user/sleeperuser",
		},
		{
			"string and int args",
			endpointLeagueMatchups,
			[]any{"123", 4},
			"https://api.sleeper.app/v1/league/123/matchups/4",
		},
		{
			"typed string args",
			endpointUserLeagues,
			[]any{"42", SportNFL, "2024"},
			"https://api.sleeper.app/v1/user/42/leagues/nfl/2024",
		},
	}

	client := &Client{baseURL: endpointBaseURL + "/v1"}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			endpoint := client.buildEndpoint(tc.template, tc.args...)
			if endpoint != tc.expected {
				t.Errorf("expected endpoint %q, got %q", tc.expected, endpoint)
				return
			}
		})
	}
}

func TestUnitEndpoint_AddQueryParams(t *testing.T) {
	tt := []struct {
		testcase string
		endpoint string
		params   map[string]string
		expected string
	}{
		{
			"nil params returns endpoint unchanged",
			"https://example.com/path",
			nil,
			"https://example.com/path",
		},
		{
			"empty params returns endpoint unchanged",
			"https://example.com/path",
			map[string]string{},
			"https://example.com/path",
		},
		{
			"single param",
			"https://example.com/path",
			map[string]string{"limit": "10"},
			"https://example.com/path?limit=10",
		},
		{
			"multiple params are sorted by key",
			"https://example.com/path",
			map[string]string{"limit": "10", "lookback_hours": "24"},
			"https://example.com/path?limit=10&lookback_hours=24",
		},
		{
			"empty values are skipped",
			"https://example.com/path",
			map[string]string{"limit": "10", "idp": ""},
			"https://example.com/path?limit=10",
		},
		{
			"only empty values leaves no query",
			"https://example.com/path",
			map[string]string{"idp": ""},
			"https://example.com/path",
		},
		{
			"existing query is kept",
			"https://example.com/path?a=1",
			map[string]string{"b": "2"},
			"https://example.com/path?a=1&b=2",
		},
		{
			"keys and values are escaped",
			"https://example.com/path",
			map[string]string{"position[]": "QB & RB"},
			"https://example.com/path?position%5B%5D=QB+%26+RB",
		},
		{
			"unparseable endpoint returns endpoint unchanged",
			"http://[::1",
			map[string]string{"limit": "10"},
			"http://[::1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			endpoint := addQueryParams(tc.endpoint, tc.params)
			if endpoint != tc.expected {
				t.Errorf("expected endpoint %q, got %q", tc.expected, endpoint)
				return
			}
		})
	}
}
