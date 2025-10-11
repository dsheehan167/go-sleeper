package sleeper

import (
	"context"
	"testing"
)

func TestSportState_Get(t *testing.T) {
	tt := []struct {
		testcase   string
		sport      sport
		shouldPass bool
	}{
		{
			"valid sport - nfl",
			SportNFL,
			true,
		},
		{
			"valid sport - nba",
			SportNBA,
			true,
		},
		{
			"valid sport - nhl",
			SportNHL,
			true,
		},
		{
			"valid sport - mlb",
			SportMLB,
			true,
		},
		{
			"invalid sport",
			"invalid_sport",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			sportState, err := testClient.GetSportState(context.Background(), tc.sport)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got sport state: %v", sportState)
				return
			}

			t.Logf("retrieved sport state: %+v", sportState)
		})
	}
}
