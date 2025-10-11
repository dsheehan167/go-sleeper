package sleeper

import (
	"context"
	"testing"
)

func TestUser_Get(t *testing.T) {
	tt := []struct {
		testcase            string
		identity            string
		expectedDisplayName string
		shouldPass          bool
	}{
		{
			"valid user_id",
			"483459259485384704",
			"SleeperUser",
			true,
		},
		{
			"valid username",
			"sleeperuser",
			"SleeperUser",
			true,
		},
		{
			"empty userID",
			"",
			"",
			false,
		},
		{
			"invalid userID",
			"invalid_user_id_483459259485384704",
			"",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			user, err := testClient.GetUser(context.Background(), tc.identity)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
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

			t.Logf("retrieved user: %+v", user)
		})
	}
}

func TestUser_GetLeagues(t *testing.T) {
	tt := []struct {
		testcase   string
		userID     string
		season     string
		sport      sport
		shouldPass bool
	}{
		{
			"valid userID with leagues",
			"476735150112768",
			"2019",
			SportNFL,
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
			"2019",
			SportNFL,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			leagues, err := testClient.GetUserLeagues(context.Background(), tc.userID, tc.sport, tc.season)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d leagues", len(leagues))
				return
			}

			if len(leagues) == 0 {
				t.Errorf("expected at least one league for user %s in season %s", tc.userID, tc.season)
				return
			}

			t.Logf("retrieved %d leagues for user %s in season %s", len(leagues), tc.userID, tc.season)
		})
	}
}
