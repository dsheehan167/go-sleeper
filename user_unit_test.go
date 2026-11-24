package sleeper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestUnitUser_Get(t *testing.T) {
	tt := []struct {
		testcase            string
		identity            string
		wantPath            string
		status              int
		body                string
		expectedDisplayName string
		wantErr             string
		shouldPass          bool
	}{
		{
			"valid user_id",
			"483459259485384704",
			"/v1/user/483459259485384704",
			http.StatusOK,
			`{"user_id":"483459259485384704","username":"sleeperuser","display_name":"SleeperUser"}`,
			"SleeperUser",
			"",
			true,
		},
		{
			"identity is trimmed",
			"  sleeperuser\t",
			"/v1/user/sleeperuser",
			http.StatusOK,
			`{"user_id":"1","username":"sleeperuser","display_name":"SleeperUser"}`,
			"SleeperUser",
			"",
			true,
		},
		{
			"empty identity",
			"",
			"",
			0,
			"",
			"",
			"userID is required",
			false,
		},
		{
			"whitespace identity",
			"   ",
			"",
			0,
			"",
			"",
			"userID is required",
			false,
		},
		{
			"not found",
			"missing",
			"/v1/user/missing",
			http.StatusNotFound,
			``,
			"",
			"getting user: API error: 404",
			false,
		},
		{
			"server error",
			"sleeperuser",
			"/v1/user/sleeperuser",
			http.StatusInternalServerError,
			``,
			"",
			"getting user: API error: 500",
			false,
		},
		{
			"bad json",
			"sleeperuser",
			"/v1/user/sleeperuser",
			http.StatusOK,
			`{"user_id":`,
			"",
			"unmarshaling user",
			false,
		},
		{
			"wrong json type",
			"sleeperuser",
			"/v1/user/sleeperuser",
			http.StatusOK,
			`{"is_bot":"yes"}`,
			"",
			"unmarshaling user",
			false,
		},
		{
			"null body",
			"sleeperuser",
			"/v1/user/sleeperuser",
			http.StatusOK,
			`null`,
			"",
			"user not found",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			user, err := client.GetUser(context.Background(), tc.identity)
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
				t.Errorf("expected failure but got user: %v", user)
				return
			}

			if user.DisplayName != tc.expectedDisplayName {
				t.Errorf("expected display name %q, got %q", tc.expectedDisplayName, user.DisplayName)
				return
			}
		})
	}
}

func TestUnitUser_GetLeagues(t *testing.T) {
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
			"476735150112768",
			SportNFL,
			"2019",
			"/v1/user/476735150112768/leagues/nfl/2019",
			http.StatusOK,
			`[{"league_id":"1","name":"One"},{"league_id":"2","name":"Two"}]`,
			2,
			"",
			true,
		},
		{
			"user id is trimmed",
			" 42 ",
			SportNBA,
			"2024",
			"/v1/user/42/leagues/nba/2024",
			http.StatusOK,
			`[{"league_id":"1"}]`,
			1,
			"",
			true,
		},
		{
			"no leagues",
			"42",
			SportNFL,
			"2024",
			"/v1/user/42/leagues/nfl/2024",
			http.StatusOK,
			`[]`,
			0,
			"",
			true,
		},
		{
			"missing fields",
			"",
			"",
			"",
			"",
			0,
			"",
			0,
			"userID is required\nsport is required\nseason is required",
			false,
		},
		{
			"missing sport",
			"42",
			"",
			"2024",
			"",
			0,
			"",
			0,
			"sport is required",
			false,
		},
		{
			"missing season",
			"42",
			SportNFL,
			"",
			"",
			0,
			"",
			0,
			"season is required",
			false,
		},
		{
			"server error",
			"42",
			SportNFL,
			"2024",
			"/v1/user/42/leagues/nfl/2024",
			http.StatusBadGateway,
			``,
			0,
			"getting leagues for user: API error: 502",
			false,
		},
		{
			"bad json",
			"42",
			SportNFL,
			"2024",
			"/v1/user/42/leagues/nfl/2024",
			http.StatusOK,
			`{"league_id":"1"}`,
			0,
			"unmarshaling leagues",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respond(t, tc.wantPath, tc.status, tc.body))

			leagues, err := client.GetUserLeagues(context.Background(), tc.userID, tc.sport, tc.season)
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
				t.Errorf("expected failure but got %d leagues", len(leagues))
				return
			}

			if len(leagues) != tc.expectedCount {
				t.Errorf("expected %d leagues, got %d", tc.expectedCount, len(leagues))
				return
			}
		})
	}
}
