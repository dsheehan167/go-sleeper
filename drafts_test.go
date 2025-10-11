package sleeper

import (
	"context"
	"testing"
)

func TestDraft_Get(t *testing.T) {
	tt := []struct {
		testcase       string
		draftID        string
		expectedRounds int
		shouldPass     bool
	}{
		{
			testcase:       "valid draft ID",
			draftID:        "483459723828391936",
			expectedRounds: 15,
			shouldPass:     true,
		},
		{
			testcase:       "missing draft ID",
			draftID:        "",
			expectedRounds: 0,
			shouldPass:     false,
		},
		{
			testcase:       "invalid draft ID",
			draftID:        "invalid",
			expectedRounds: 0,
			shouldPass:     false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			draft, err := testClient.GetDraft(context.Background(), tc.draftID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got draft: %v", draft)
				return
			}

			if draft.Settings == nil {
				t.Errorf("expected draft settings but got nil")
				return
			}

			if draft.Settings.Rounds != tc.expectedRounds {
				t.Errorf("expected %d rounds, got %d", tc.expectedRounds, draft.Settings.Rounds)
				return
			}

			t.Logf("retrieved draft with %d rounds", draft.Settings.Rounds)
		})
	}
}

func TestDraft_UserDrafts(t *testing.T) {
	tt := []struct {
		testcase   string
		userID     string
		sport      sport
		season     string
		shouldPass bool
	}{
		{
			"valid userID with drafts",
			"476735150112768",
			SportNFL,
			"2019",
			true,
		},
		{
			"missing fields",
			"",
			"",
			"",
			false,
		},
		{
			"invalid userID",
			"invalid",
			SportNFL,
			"2019",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			drafts, err := testClient.GetUserDrafts(context.Background(), tc.userID, tc.sport, tc.season)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d drafts", len(drafts))
				return
			}

			if len(drafts) == 0 {
				t.Errorf("expected at least one draft for user %s in season %s", tc.userID, tc.season)
				return
			}

			t.Logf("retrieved %d drafts", len(drafts))
		})
	}
}

func TestDraft_LeagueDrafts(t *testing.T) {
	tt := []struct {
		testcase   string
		leagueID   string
		shouldPass bool
	}{
		{
			"valid league ID",
			"1257116530085208064",
			true,
		},
		{
			"missing league ID",
			"",
			false,
		},
		{
			"invalid league ID",
			"invalid",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			drafts, err := testClient.GetLeagueDrafts(context.Background(), tc.leagueID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d drafts", len(drafts))
				return
			}

			if len(drafts) == 0 {
				t.Errorf("expected at least one draft for league %s", tc.leagueID)
				return
			}

			t.Logf("retrieved %d drafts", len(drafts))
		})
	}
}

func TestDraft_Picks(t *testing.T) {
	tt := []struct {
		testcase       string
		draftID        string
		expectedPickCt int
		shouldPass     bool
	}{
		{
			"valid draft ID",
			"289646328508579840",
			180,
			true,
		},
		{
			"missing draft ID",
			"",
			0,
			false,
		},
		{
			"invalid draft ID",
			"invalid",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			picks, err := testClient.GetDraftPicks(context.Background(), tc.draftID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d picks", len(picks))
				return
			}

			if len(picks) != tc.expectedPickCt {
				t.Errorf("expected %d picks, got %d", tc.expectedPickCt, len(picks))
				return
			}

			t.Logf("retrieved %d picks", len(picks))
		})
	}
}

func TestDraft_TradedPicks(t *testing.T) {
	tt := []struct {
		testcase       string
		draftID        string
		expectedPickCt int
		shouldPass     bool
	}{
		{
			"valid draft ID",
			"1234907380219641856",
			14,
			true,
		},
		{
			"missing draft ID",
			"",
			0,
			false,
		},
		{
			"invalid draft ID",
			"invalid",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			picks, err := testClient.GetDraftTradedPicks(context.Background(), tc.draftID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d picks", len(picks))
				return
			}

			if len(picks) != tc.expectedPickCt {
				t.Errorf("expected %d picks, got %d", tc.expectedPickCt, len(picks))
				return
			}

			t.Logf("retrieved %d picks", len(picks))
		})
	}
}
