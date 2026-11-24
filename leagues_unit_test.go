package sleeper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestUnitLeague_Get(t *testing.T) {
	tt := []struct {
		testcase     string
		leagueID     string
		wantPath     string
		status       int
		body         string
		expectedName string
		wantErr      string
		shouldPass   bool
	}{
		{
			"valid leagueID",
			"289646328504385536",
			"/v1/league/289646328504385536",
			http.StatusOK,
			`{"league_id":"289646328504385536","name":"Sleeper Friends League","season":"2026","status":"in_season"}`,
			"Sleeper Friends League",
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			"/v1/league/289646328504385536",
			http.StatusOK,
			`{"league_id":"289646328504385536","name":"Trimmed"}`,
			"Trimmed",
			"",
			true,
		},
		{
			"empty leagueID",
			"",
			"",
			0,
			"",
			"",
			"leagueID is required",
			false,
		},
		{
			"not found",
			"missing",
			"/v1/league/missing",
			http.StatusNotFound,
			``,
			"",
			"getting league: API error: 404",
			false,
		},
		{
			"server error",
			"289646328504385536",
			"/v1/league/289646328504385536",
			http.StatusInternalServerError,
			``,
			"",
			"getting league: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			"/v1/league/289646328504385536",
			http.StatusOK,
			`{"name":42}`,
			"",
			"unmarshaling league",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			league, err := client.GetLeague(context.Background(), tc.leagueID)
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
				t.Errorf("expected failure but got league: %v", league)
				return
			}

			if league == nil {
				t.Errorf("expected league, got nil")
				return
			}

			if league.Name != tc.expectedName {
				t.Errorf("expected league name %q, got %q", tc.expectedName, league.Name)
				return
			}
		})
	}
}

func TestUnitLeague_GetMatchups(t *testing.T) {
	tt := []struct {
		testcase      string
		leagueID      string
		week          int
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid leagueID and week",
			"289646328504385536",
			1,
			"/v1/league/289646328504385536/matchups/1",
			http.StatusOK,
			`[{"roster_id":1,"matchup_id":1,"points":120.5},{"roster_id":2,"matchup_id":1,"points":98.2}]`,
			2,
			"",
			true,
		},
		{
			"late season week",
			"289646328504385536",
			17,
			"/v1/league/289646328504385536/matchups/17",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"missing leagueID and week",
			"",
			0,
			"",
			0,
			"",
			0,
			"leagueID is required\nweek must be greater than or equal to 1",
			false,
		},
		{
			"negative week",
			"289646328504385536",
			-1,
			"",
			0,
			"",
			0,
			"week must be greater than or equal to 1",
			false,
		},
		{
			"server error",
			"289646328504385536",
			1,
			"/v1/league/289646328504385536/matchups/1",
			http.StatusInternalServerError,
			``,
			0,
			"getting league matchups: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			1,
			"/v1/league/289646328504385536/matchups/1",
			http.StatusOK,
			`[{"points":"lots"}]`,
			0,
			"unmarshaling league matchups",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			matchups, err := client.GetLeagueMatchups(context.Background(), tc.leagueID, tc.week)
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
				t.Errorf("expected failure but got %d matchups", len(matchups))
				return
			}

			if len(matchups) != tc.expectedCount {
				t.Errorf("expected %d matchups, got %d", tc.expectedCount, len(matchups))
				return
			}
		})
	}
}

func TestUnitLeague_GetTransactions(t *testing.T) {
	tt := []struct {
		testcase      string
		leagueID      string
		week          int
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid leagueID and week",
			"289646328504385536",
			3,
			"/v1/league/289646328504385536/transactions/3",
			http.StatusOK,
			`[{"transaction_id":"1","type":"trade","status":"complete"},{"transaction_id":"2","type":"waiver","status":"failed"}]`,
			2,
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			3,
			"/v1/league/289646328504385536/transactions/3",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"missing leagueID",
			"",
			3,
			"",
			0,
			"",
			0,
			"leagueID is required",
			false,
		},
		{
			"zero week",
			"289646328504385536",
			0,
			"",
			0,
			"",
			0,
			"week must be greater than or equal to 1",
			false,
		},
		{
			"server error",
			"289646328504385536",
			3,
			"/v1/league/289646328504385536/transactions/3",
			http.StatusInternalServerError,
			``,
			0,
			"getting league transactions: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			3,
			"/v1/league/289646328504385536/transactions/3",
			http.StatusOK,
			`[{"transaction_id":1}]`,
			0,
			"unmarshaling league transactions",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			transactions, err := client.GetTransactions(context.Background(), tc.leagueID, tc.week)
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
				t.Errorf("expected failure but got %d transactions", len(transactions))
				return
			}

			if len(transactions) != tc.expectedCount {
				t.Errorf("expected %d transactions, got %d", tc.expectedCount, len(transactions))
				return
			}
		})
	}
}

func TestUnitLeague_GetRosters(t *testing.T) {
	tt := []struct {
		testcase      string
		leagueID      string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid leagueID",
			"289646328504385536",
			"/v1/league/289646328504385536/rosters",
			http.StatusOK,
			`[{"roster_id":1,"owner_id":"9","players":["4034"],"metadata":{"record":"W","p_nick_4034":"CMC"}},{"roster_id":2}]`,
			2,
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			"/v1/league/289646328504385536/rosters",
			http.StatusOK,
			`[{"roster_id":1,"owner_id":"9","players":["4034"],"metadata":{"record":"W","p_nick_4034":"CMC"}},{"roster_id":2}]`,
			2,
			"",
			true,
		},
		{
			"empty list",
			"289646328504385536",
			"/v1/league/289646328504385536/rosters",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"empty leagueID",
			"",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"whitespace leagueID",
			"   ",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"not found",
			"289646328504385536",
			"/v1/league/289646328504385536/rosters",
			http.StatusNotFound,
			``,
			0,
			"getting league rosters: API error: 404",
			false,
		},
		{
			"server error",
			"289646328504385536",
			"/v1/league/289646328504385536/rosters",
			http.StatusInternalServerError,
			``,
			0,
			"getting league rosters: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			"/v1/league/289646328504385536/rosters",
			http.StatusOK,
			`[{"roster_id":"one"}]`,
			0,
			"unmarshaling league rosters",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetLeagueRosters(context.Background(), tc.leagueID)
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
				t.Errorf("expected failure but got %d results", len(results))
				return
			}

			if len(results) != tc.expectedCount {
				t.Errorf("expected %d results, got %d", tc.expectedCount, len(results))
				return
			}
		})
	}
}

func TestUnitLeague_GetUsers(t *testing.T) {
	tt := []struct {
		testcase      string
		leagueID      string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid leagueID",
			"289646328504385536",
			"/v1/league/289646328504385536/users",
			http.StatusOK,
			`[{"user_id":"1","display_name":"One"},{"user_id":"2","display_name":"Two"},{"user_id":"3","is_bot":true}]`,
			3,
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			"/v1/league/289646328504385536/users",
			http.StatusOK,
			`[{"user_id":"1","display_name":"One"},{"user_id":"2","display_name":"Two"},{"user_id":"3","is_bot":true}]`,
			3,
			"",
			true,
		},
		{
			"empty list",
			"289646328504385536",
			"/v1/league/289646328504385536/users",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"empty leagueID",
			"",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"whitespace leagueID",
			"   ",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"not found",
			"289646328504385536",
			"/v1/league/289646328504385536/users",
			http.StatusNotFound,
			``,
			0,
			"getting league users: API error: 404",
			false,
		},
		{
			"server error",
			"289646328504385536",
			"/v1/league/289646328504385536/users",
			http.StatusInternalServerError,
			``,
			0,
			"getting league users: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			"/v1/league/289646328504385536/users",
			http.StatusOK,
			`{"user_id":"1"}`,
			0,
			"unmarshaling league users",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetLeagueUsers(context.Background(), tc.leagueID)
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
				t.Errorf("expected failure but got %d results", len(results))
				return
			}

			if len(results) != tc.expectedCount {
				t.Errorf("expected %d results, got %d", tc.expectedCount, len(results))
				return
			}
		})
	}
}

func TestUnitLeague_GetTradedPicks(t *testing.T) {
	tt := []struct {
		testcase      string
		leagueID      string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid leagueID",
			"289646328504385536",
			"/v1/league/289646328504385536/traded_picks",
			http.StatusOK,
			`[{"season":"2027","round":1,"roster_id":3,"previous_owner_id":3,"owner_id":5}]`,
			1,
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			"/v1/league/289646328504385536/traded_picks",
			http.StatusOK,
			`[{"season":"2027","round":1,"roster_id":3,"previous_owner_id":3,"owner_id":5}]`,
			1,
			"",
			true,
		},
		{
			"empty list",
			"289646328504385536",
			"/v1/league/289646328504385536/traded_picks",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"empty leagueID",
			"",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"whitespace leagueID",
			"   ",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"not found",
			"289646328504385536",
			"/v1/league/289646328504385536/traded_picks",
			http.StatusNotFound,
			``,
			0,
			"getting traded picks: API error: 404",
			false,
		},
		{
			"server error",
			"289646328504385536",
			"/v1/league/289646328504385536/traded_picks",
			http.StatusInternalServerError,
			``,
			0,
			"getting traded picks: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			"/v1/league/289646328504385536/traded_picks",
			http.StatusOK,
			`[{"round":"first"}]`,
			0,
			"unmarshaling traded picks",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetLeagueTradedPicks(context.Background(), tc.leagueID)
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
				t.Errorf("expected failure but got %d results", len(results))
				return
			}

			if len(results) != tc.expectedCount {
				t.Errorf("expected %d results, got %d", tc.expectedCount, len(results))
				return
			}
		})
	}
}

func TestUnitLeague_GetWinnersBracket(t *testing.T) {
	tt := []struct {
		testcase      string
		leagueID      string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid leagueID",
			"289646328504385536",
			"/v1/league/289646328504385536/winners_bracket",
			http.StatusOK,
			`[{"r":1,"m":1,"t1":1,"t2":4,"w":1,"l":4},{"r":2,"m":3,"t1_from":{"w":1},"t2_from":{"w":2}}]`,
			2,
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			"/v1/league/289646328504385536/winners_bracket",
			http.StatusOK,
			`[{"r":1,"m":1,"t1":1,"t2":4,"w":1,"l":4},{"r":2,"m":3,"t1_from":{"w":1},"t2_from":{"w":2}}]`,
			2,
			"",
			true,
		},
		{
			"empty list",
			"289646328504385536",
			"/v1/league/289646328504385536/winners_bracket",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"empty leagueID",
			"",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"whitespace leagueID",
			"   ",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"not found",
			"289646328504385536",
			"/v1/league/289646328504385536/winners_bracket",
			http.StatusNotFound,
			``,
			0,
			"getting winners bracket: API error: 404",
			false,
		},
		{
			"server error",
			"289646328504385536",
			"/v1/league/289646328504385536/winners_bracket",
			http.StatusInternalServerError,
			``,
			0,
			"getting winners bracket: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			"/v1/league/289646328504385536/winners_bracket",
			http.StatusOK,
			`[{"r":1`,
			0,
			"unmarshaling winners bracket",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetLeagueWinnersBracket(context.Background(), tc.leagueID)
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
				t.Errorf("expected failure but got %d results", len(results))
				return
			}

			if len(results) != tc.expectedCount {
				t.Errorf("expected %d results, got %d", tc.expectedCount, len(results))
				return
			}
		})
	}
}

func TestUnitLeague_GetLosersBracket(t *testing.T) {
	tt := []struct {
		testcase      string
		leagueID      string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid leagueID",
			"289646328504385536",
			"/v1/league/289646328504385536/losers_bracket",
			http.StatusOK,
			`[{"r":1,"m":1,"t1":5,"t2":6,"p":7}]`,
			1,
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			"/v1/league/289646328504385536/losers_bracket",
			http.StatusOK,
			`[{"r":1,"m":1,"t1":5,"t2":6,"p":7}]`,
			1,
			"",
			true,
		},
		{
			"empty list",
			"289646328504385536",
			"/v1/league/289646328504385536/losers_bracket",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"empty leagueID",
			"",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"whitespace leagueID",
			"   ",
			"",
			0,
			``,
			0,
			"leagueID is required",
			false,
		},
		{
			"not found",
			"289646328504385536",
			"/v1/league/289646328504385536/losers_bracket",
			http.StatusNotFound,
			``,
			0,
			"getting losers bracket: API error: 404",
			false,
		},
		{
			"server error",
			"289646328504385536",
			"/v1/league/289646328504385536/losers_bracket",
			http.StatusInternalServerError,
			``,
			0,
			"getting losers bracket: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			"/v1/league/289646328504385536/losers_bracket",
			http.StatusOK,
			`[{"t1":"five"}]`,
			0,
			"unmarshaling losers bracket",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetLeagueLosersBracket(context.Background(), tc.leagueID)
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
				t.Errorf("expected failure but got %d results", len(results))
				return
			}

			if len(results) != tc.expectedCount {
				t.Errorf("expected %d results, got %d", tc.expectedCount, len(results))
				return
			}
		})
	}
}
