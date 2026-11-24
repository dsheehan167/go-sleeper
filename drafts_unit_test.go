package sleeper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestUnitDraft_Get(t *testing.T) {
	tt := []struct {
		testcase       string
		draftID        string
		wantPath       string
		status         int
		body           string
		expectedStatus string
		wantErr        string
		shouldPass     bool
	}{
		{
			"valid draftID",
			"257270637750382592",
			"/v1/draft/257270637750382592",
			http.StatusOK,
			`{"draft_id":"257270637750382592","league_id":"1","season":"2026","status":"complete","type":"snake"}`,
			"complete",
			"",
			true,
		},
		{
			"draftID is trimmed",
			" 257270637750382592\n",
			"/v1/draft/257270637750382592",
			http.StatusOK,
			`{"draft_id":"257270637750382592","status":"pre_draft"}`,
			"pre_draft",
			"",
			true,
		},
		{
			"empty draftID",
			"",
			"",
			0,
			"",
			"",
			"draftID is required",
			false,
		},
		{
			"not found",
			"missing",
			"/v1/draft/missing",
			http.StatusNotFound,
			``,
			"",
			"getting draft: API error: 404",
			false,
		},
		{
			"bad json",
			"257270637750382592",
			"/v1/draft/257270637750382592",
			http.StatusOK,
			`{"draft_id":257270637750382592}`,
			"",
			"unmarshaling draft",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			draft, err := client.GetDraft(context.Background(), tc.draftID)
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
				t.Errorf("expected failure but got draft: %v", draft)
				return
			}

			if draft == nil {
				t.Errorf("expected draft, got nil")
				return
			}

			if draft.Status != tc.expectedStatus {
				t.Errorf("expected draft status %q, got %q", tc.expectedStatus, draft.Status)
				return
			}
		})
	}
}

func TestUnitDraft_UserDrafts(t *testing.T) {
	tt := []struct {
		testcase      string
		userID        string
		sport         Sport
		season        string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid request",
			"457511950237696",
			SportNFL,
			"2026",
			"/v1/user/457511950237696/drafts/nfl/2026",
			http.StatusOK,
			`[{"draft_id":"1"},{"draft_id":"2"}]`,
			2,
			"",
			true,
		},
		{
			"userID and season are trimmed",
			" 457511950237696 ",
			SportNFL,
			" 2026 ",
			"/v1/user/457511950237696/drafts/nfl/2026",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"missing fields",
			" ",
			"",
			" ",
			"",
			0,
			"",
			0,
			"userID is required\nsport is required\nseason is required",
			false,
		},
		{
			"missing sport",
			"457511950237696",
			"",
			"2026",
			"",
			0,
			"",
			0,
			"sport is required",
			false,
		},
		{
			"server error",
			"457511950237696",
			SportNFL,
			"2026",
			"/v1/user/457511950237696/drafts/nfl/2026",
			http.StatusInternalServerError,
			``,
			0,
			"getting drafts for user: API error: 500",
			false,
		},
		{
			"bad json",
			"457511950237696",
			SportNFL,
			"2026",
			"/v1/user/457511950237696/drafts/nfl/2026",
			http.StatusOK,
			`{"draft_id":"1"}`,
			0,
			"unmarshaling drafts",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			drafts, err := client.GetUserDrafts(context.Background(), tc.userID, tc.sport, tc.season)
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
				t.Errorf("expected failure but got %d drafts", len(drafts))
				return
			}

			if len(drafts) != tc.expectedCount {
				t.Errorf("expected %d drafts, got %d", tc.expectedCount, len(drafts))
				return
			}
		})
	}
}

func TestUnitDraft_LeagueDrafts(t *testing.T) {
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
			"/v1/league/289646328504385536/drafts",
			http.StatusOK,
			`[{"draft_id":"1","league_id":"289646328504385536","status":"complete"}]`,
			1,
			"",
			true,
		},
		{
			"leagueID is trimmed",
			" 289646328504385536 ",
			"/v1/league/289646328504385536/drafts",
			http.StatusOK,
			`[{"draft_id":"1","league_id":"289646328504385536","status":"complete"}]`,
			1,
			"",
			true,
		},
		{
			"empty list",
			"289646328504385536",
			"/v1/league/289646328504385536/drafts",
			http.StatusOK,
			`[]`,
			0,
			"no drafts found for league",
			false,
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
			"/v1/league/289646328504385536/drafts",
			http.StatusNotFound,
			``,
			0,
			"getting drafts for league: API error: 404",
			false,
		},
		{
			"server error",
			"289646328504385536",
			"/v1/league/289646328504385536/drafts",
			http.StatusInternalServerError,
			``,
			0,
			"getting drafts for league: API error: 500",
			false,
		},
		{
			"bad json",
			"289646328504385536",
			"/v1/league/289646328504385536/drafts",
			http.StatusOK,
			`{"draft_id":"1"}`,
			0,
			"unmarshaling drafts",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetLeagueDrafts(context.Background(), tc.leagueID)
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

func TestUnitDraft_TradedPicks(t *testing.T) {
	tt := []struct {
		testcase      string
		draftID       string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid draftID",
			"257270637750382592",
			"/v1/draft/257270637750382592/traded_picks",
			http.StatusOK,
			`[{"season":"2026","round":2,"roster_id":1,"previous_owner_id":1,"owner_id":2}]`,
			1,
			"",
			true,
		},
		{
			"draftID is trimmed",
			" 257270637750382592 ",
			"/v1/draft/257270637750382592/traded_picks",
			http.StatusOK,
			`[{"season":"2026","round":2,"roster_id":1,"previous_owner_id":1,"owner_id":2}]`,
			1,
			"",
			true,
		},
		{
			"empty list",
			"257270637750382592",
			"/v1/draft/257270637750382592/traded_picks",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"empty draftID",
			"",
			"",
			0,
			``,
			0,
			"draftID is required",
			false,
		},
		{
			"whitespace draftID",
			"   ",
			"",
			0,
			``,
			0,
			"draftID is required",
			false,
		},
		{
			"not found",
			"257270637750382592",
			"/v1/draft/257270637750382592/traded_picks",
			http.StatusNotFound,
			``,
			0,
			"getting traded picks for draft: API error: 404",
			false,
		},
		{
			"server error",
			"257270637750382592",
			"/v1/draft/257270637750382592/traded_picks",
			http.StatusInternalServerError,
			``,
			0,
			"getting traded picks for draft: API error: 500",
			false,
		},
		{
			"bad json",
			"257270637750382592",
			"/v1/draft/257270637750382592/traded_picks",
			http.StatusOK,
			`[{"owner_id":"two"}]`,
			0,
			"unmarshaling traded picks",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetDraftTradedPicks(context.Background(), tc.draftID)
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

func TestUnitDraft_Picks(t *testing.T) {
	tt := []struct {
		testcase      string
		draftID       string
		wantPath      string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid draftID",
			"257270637750382592",
			"/v1/draft/257270637750382592/picks",
			http.StatusOK,
			`[{"draft_id":"257270637750382592","pick_no":1,"round":1,"player_id":"4034","roster_id":1},{"pick_no":2,"round":1,"player_id":"6794","roster_id":2}]`,
			2,
			"",
			true,
		},
		{
			"draftID is trimmed",
			" 257270637750382592 ",
			"/v1/draft/257270637750382592/picks",
			http.StatusOK,
			`[{"draft_id":"257270637750382592","pick_no":1,"round":1,"player_id":"4034","roster_id":1},{"pick_no":2,"round":1,"player_id":"6794","roster_id":2}]`,
			2,
			"",
			true,
		},
		{
			"empty list",
			"257270637750382592",
			"/v1/draft/257270637750382592/picks",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"empty draftID",
			"",
			"",
			0,
			``,
			0,
			"draftID is required",
			false,
		},
		{
			"whitespace draftID",
			"   ",
			"",
			0,
			``,
			0,
			"draftID is required",
			false,
		},
		{
			"not found",
			"257270637750382592",
			"/v1/draft/257270637750382592/picks",
			http.StatusNotFound,
			``,
			0,
			"getting picks for draft: API error: 404",
			false,
		},
		{
			"server error",
			"257270637750382592",
			"/v1/draft/257270637750382592/picks",
			http.StatusInternalServerError,
			``,
			0,
			"getting picks for draft: API error: 500",
			false,
		},
		{
			"bad json",
			"257270637750382592",
			"/v1/draft/257270637750382592/picks",
			http.StatusOK,
			`[{"pick_no":"one"}]`,
			0,
			"unmarshaling picks",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			results, err := client.GetDraftPicks(context.Background(), tc.draftID)
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
