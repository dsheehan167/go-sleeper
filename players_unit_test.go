package sleeper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestUnitPlayer_List(t *testing.T) {
	tt := []struct {
		testcase      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid players",
			http.StatusOK,
			`{"4034":{"player_id":"4034","full_name":"Christian McCaffrey","position":"RB","team":"SF"},"CAR":{"player_id":"CAR","position":"DEF","team":null}}`,
			2,
			"",
			true,
		},
		{
			"empty map",
			http.StatusOK,
			`{}`,
			0,
			"",
			true,
		},
		{
			"server error",
			http.StatusInternalServerError,
			``,
			0,
			"getting players: API error: 500",
			false,
		},
		{
			"bad json",
			http.StatusOK,
			`[{"player_id":"4034"}]`,
			0,
			"unmarshaling players",
			false,
		},
		{
			"null body",
			http.StatusOK,
			`null`,
			0,
			"players not found",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, "/v1/players/nfl", tc.status, tc.body))

			players, err := client.ListNFLPlayers(context.Background())
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
				t.Errorf("expected failure but got %d players", len(players))
				return
			}

			if len(players) != tc.expectedCount {
				t.Errorf("expected %d players, got %d", tc.expectedCount, len(players))
				return
			}
		})
	}
}

func TestUnitPlayer_ListTrending(t *testing.T) {
	tt := []struct {
		testcase      string
		sport         Sport
		trendingType  TrendingType
		options       TrendingPlayerOptions
		wantPath      string
		wantQuery     string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid add without options",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{},
			"/v1/players/nfl/trending/add",
			"",
			http.StatusOK,
			`[{"player_id":"4034","count":120},{"player_id":"6794","count":80}]`,
			2,
			"",
			true,
		},
		{
			"valid drop with options",
			SportNBA,
			TrendingTypeDrop,
			TrendingPlayerOptions{LookbackHours: 48, Limit: 25},
			"/v1/players/nba/trending/drop",
			"limit=25&lookback_hours=48",
			http.StatusOK,
			`[{"player_id":"1"}]`,
			1,
			"",
			true,
		},
		{
			"only limit",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{Limit: 50},
			"/v1/players/nfl/trending/add",
			"limit=50",
			http.StatusOK,
			`[{"player_id":"1"}]`,
			1,
			"",
			true,
		},
		{
			"missing sport and type",
			"",
			"",
			TrendingPlayerOptions{},
			"",
			"",
			0,
			"",
			0,
			"sport is required\ntrendingType is required",
			false,
		},
		{
			"negative lookback hours",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{LookbackHours: -1},
			"",
			"",
			0,
			"",
			0,
			"LookbackHours must be greater than or equal to zero",
			false,
		},
		{
			"negative limit",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{Limit: -1},
			"",
			"",
			0,
			"",
			0,
			"Limit must be greater than or equal to zero",
			false,
		},
		{
			"limit above maximum",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{Limit: 51},
			"",
			"",
			0,
			"",
			0,
			"Limit cannot exceed 50",
			false,
		},
		{
			"server error",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{},
			"/v1/players/nfl/trending/add",
			"",
			http.StatusInternalServerError,
			``,
			0,
			"getting trending players: API error: 500",
			false,
		},
		{
			"bad json",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{},
			"/v1/players/nfl/trending/add",
			"",
			http.StatusOK,
			`{"player_id":"1"}`,
			0,
			"unmarshaling trending players",
			false,
		},
		{
			"empty list",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{},
			"/v1/players/nfl/trending/add",
			"",
			http.StatusOK,
			`[]`,
			0,
			"trending players not found",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respondExact(t, "", tc.wantPath, tc.wantQuery, tc.status, tc.body))

			players, err := client.ListTrendingPlayers(context.Background(), tc.sport, tc.trendingType, tc.options)
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
				t.Errorf("expected failure but got %d players", len(players))
				return
			}

			if len(players) != tc.expectedCount {
				t.Errorf("expected %d players, got %d", tc.expectedCount, len(players))
				return
			}
		})
	}
}
