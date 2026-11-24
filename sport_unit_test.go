package sleeper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestUnitSportState_Get(t *testing.T) {
	tt := []struct {
		testcase         string
		sport            Sport
		wantPath         string
		status           int
		body             string
		expectedWeek     int
		expectedPrevious FlexibleString
		wantErr          string
		shouldPass       bool
	}{
		{
			"valid nfl state",
			SportNFL,
			"/v1/state/nfl",
			http.StatusOK,
			`{"week":12,"season_type":"regular","season":"2026","previous_season":"2025","display_week":12}`,
			12,
			"2025",
			"",
			true,
		},
		{
			"previous season as int",
			SportNBA,
			"/v1/state/nba",
			http.StatusOK,
			`{"week":4,"season":"2026","previous_season":2025}`,
			4,
			"2025",
			"",
			true,
		},
		{
			"empty sport",
			"",
			"",
			0,
			"",
			0,
			"",
			"sport is required",
			false,
		},
		{
			"server error",
			SportNFL,
			"/v1/state/nfl",
			http.StatusServiceUnavailable,
			``,
			0,
			"",
			"getting sport state: API error: 503",
			false,
		},
		{
			"bad json",
			SportNFL,
			"/v1/state/nfl",
			http.StatusOK,
			`{"week":"twelve"}`,
			0,
			"",
			"unmarshaling sport state",
			false,
		},
		{
			"invalid previous season",
			SportNFL,
			"/v1/state/nfl",
			http.StatusOK,
			`{"previous_season":true}`,
			0,
			"",
			"unmarshaling sport state",
			false,
		},
		{
			"null body",
			SportNFL,
			"/v1/state/nfl",
			http.StatusOK,
			`null`,
			0,
			"",
			"sport state not found",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			state, err := client.GetSportState(context.Background(), tc.sport)
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
				t.Errorf("expected failure but got state: %+v", state)
				return
			}

			if state.Week != tc.expectedWeek {
				t.Errorf("expected week %d, got %d", tc.expectedWeek, state.Week)
			}
			if state.PreviousSeason != tc.expectedPrevious {
				t.Errorf("expected previous season %q, got %q", tc.expectedPrevious, state.PreviousSeason)
			}
		})
	}
}
