package sleeper

import (
	"context"
	"testing"
)

func TestLeague_Get(t *testing.T) {
	tt := []struct {
		testcase     string
		leagueID     string
		expectedName string
		shouldPass   bool
	}{
		{
			"valid league ID",
			"1257116530085208064",
			"Family League",
			true,
		},
		{
			"missing league ID",
			"",
			"",
			false,
		},
		{
			"invalid league ID",
			"invalid",
			"",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			league, err := testClient.GetLeague(context.Background(), tc.leagueID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got league: %v", league)
				return
			}

			if league.Name != tc.expectedName {
				t.Errorf("expected league name %q, got %q", tc.expectedName, league.Name)
				return
			}

			t.Logf("retrieved league: %+v", league)
		})
	}
}

func TestLeague_GetRosters(t *testing.T) {
	tt := []struct {
		testcase         string
		leagueID         string
		expectedRosterCt int
		shouldPass       bool
	}{
		{
			"valid league ID",
			"289646328504385536",
			12,
			true,
		},
		{
			"missing league ID",
			"",
			0,
			false,
		},
		{
			"invalid league ID",
			"invalid",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			rosters, err := testClient.GetLeagueRosters(context.Background(), tc.leagueID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got rosters: %v", rosters)
				return
			}

			if len(rosters) != tc.expectedRosterCt {
				t.Errorf("expected %d rosters, got %d", tc.expectedRosterCt, len(rosters))
				return
			}

			t.Logf("retrieved %d rosters", len(rosters))
		})
	}
}

func TestLeague_GetUsers(t *testing.T) {
	tt := []struct {
		testcase       string
		leagueID       string
		expectedUserCt int
		shouldPass     bool
	}{
		{
			"valid league ID",
			"289646328504385536",
			14,
			true,
		},
		{
			"missing league ID",
			"",
			0,
			false,
		},
		{
			"invalid league ID",
			"invalid",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			users, err := testClient.GetLeagueUsers(context.Background(), tc.leagueID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got users: %v", users)
				return
			}

			if len(users) != tc.expectedUserCt {
				t.Errorf("expected %d users, got %d", tc.expectedUserCt, len(users))
				return
			}

			t.Logf("retrieved %d users", len(users))
		})
	}
}

func TestLeague_GetMatchups(t *testing.T) {
	tt := []struct {
		testcase          string
		leagueID          string
		week              int
		expectedMatchupCt int
		shouldPass        bool
	}{
		{
			"valid league ID and week",
			"289646328504385536",
			1,
			12,
			true,
		},
		{
			"week 0",
			"289646328504385536",
			0,
			0,
			false,
		},
		{
			"missing league ID",
			"",
			1,
			0,
			false,
		},
		{
			"invalid league ID",
			"invalid",
			1,
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			matchups, err := testClient.GetLeagueMatchups(context.Background(), tc.leagueID, tc.week)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got matchups: %v", matchups)
				return
			}

			if len(matchups) != tc.expectedMatchupCt {
				t.Errorf("expected %d matchups, got %d", tc.expectedMatchupCt, len(matchups))
				return
			}

			t.Logf("retrieved %d matchups", len(matchups))
		})
	}
}

func TestLeague_GetTransactions(t *testing.T) {
	tt := []struct {
		testcase              string
		leagueID              string
		week                  int
		expectedTransactionCt int
		shouldPass            bool
	}{
		{
			"valid league ID",
			"289646328504385536",
			3,
			66,
			true,
		},
		{
			"week set to 0",
			"289646328504385536",
			0,
			0,
			false,
		},
		{
			"missing league ID",
			"",
			3,
			0,
			false,
		},
		{
			"invalid league ID",
			"invalid",
			3,
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			transactions, err := testClient.GetTransactions(context.Background(), tc.leagueID, tc.week)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got transactions: %v", transactions)
				return
			}

			if len(transactions) != tc.expectedTransactionCt {
				t.Errorf("expected %d transactions, got %d", tc.expectedTransactionCt, len(transactions))
				return
			}

			t.Logf("retrieved %d transactions", len(transactions))
		})
	}
}

func TestLeague_GetTradedPicks(t *testing.T) {
	tt := []struct {
		testcase       string
		leagueID       string
		expectedPickCt int
		shouldPass     bool
	}{
		{
			"valid league ID",
			"289646328504385536",
			4,
			true,
		},
		{
			"missing league ID",
			"",
			0,
			false,
		},
		{
			"invalid league ID",
			"invalid",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			picks, err := testClient.GetLeagueTradedPicks(context.Background(), tc.leagueID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got picks: %v", picks)
				return
			}

			if len(picks) != tc.expectedPickCt {
				t.Errorf("expected %d traded picks, got %d", tc.expectedPickCt, len(picks))
				return
			}

			t.Logf("retrieved %d traded picks", len(picks))
		})
	}
}

func TestLeague_GetWinnersBracket(t *testing.T) {
	tt := []struct {
		testcase          string
		leagueID          string
		expectedMatchupCt int
		shouldPass        bool
	}{
		{
			"valid league ID",
			"289646328504385536",
			7,
			true,
		},
		{
			"missing league ID",
			"",
			0,
			false,
		},
		{
			"invalid league ID",
			"invalid",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			matchups, err := testClient.GetLeagueWinnersBracket(context.Background(), tc.leagueID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got matchups: %v", matchups)
				return
			}

			if len(matchups) != tc.expectedMatchupCt {
				t.Errorf("expected %d matchups, got %d", tc.expectedMatchupCt, len(matchups))
				return
			}

			t.Logf("retrieved %d matchups", len(matchups))
		})
	}
}

func TestLeague_GetLosersBracket(t *testing.T) {
	tt := []struct {
		testcase          string
		leagueID          string
		expectedMatchupCt int
		shouldPass        bool
	}{
		{
			"valid league ID",
			"289646328504385536",
			7,
			true,
		},
		{
			"missing league ID",
			"",
			0,
			false,
		},
		{
			"invalid league ID",
			"invalid",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			matchups, err := testClient.GetLeagueLosersBracket(context.Background(), tc.leagueID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got matchups: %v", matchups)
				return
			}

			if len(matchups) != tc.expectedMatchupCt {
				t.Errorf("expected %d matchups, got %d", tc.expectedMatchupCt, len(matchups))
				return
			}

			t.Logf("retrieved %d matchups", len(matchups))
		})
	}
}
