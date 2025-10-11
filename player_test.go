package sleeper

import (
	"context"
	"testing"
)

func TestPlayer_List(t *testing.T) {
	tt := []struct {
		testcase   string
		shouldPass bool
	}{
		// {
		// 	"list players",
		// 	true,
		// },
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			players, err := testClient.ListNFLPlayers(context.Background())
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got players: %v", players)
				return
			}

			if len(players) == 0 {
				t.Errorf("expected non-empty player list")
				return
			}

			t.Logf("retrieved %d players", len(players))
		})
	}
}

func TestPlayer_ListTrending(t *testing.T) {
	tt := []struct {
		testcase     string
		sport        sport
		trendingType trendingType
		options      TrendingPlayerOptions
		shouldPass   bool
	}{
		{
			"valid sport and type - nfl - add",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{},
			true,
		},
		{
			"valid sport and type - nfl - drop",
			SportNFL,
			TrendingTypeDrop,
			TrendingPlayerOptions{},
			true,
		},
		{
			"valid sport and type - nfl - add with limit",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{Limit: 5},
			true,
		},
		{
			"valid sport and type - nfl - with lookback hours",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{LookbackHours: 12},
			true,
		},
		{
			"invalid sport",
			"invalid_sport",
			TrendingTypeAdd,
			TrendingPlayerOptions{},
			false,
		},
		{
			"invalid type",
			SportNFL,
			"invalid_type",
			TrendingPlayerOptions{},
			false,
		},
		{
			"invalid options - limit too high",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{Limit: 51},
			false,
		},
		{
			"invalid options - lookback hours negative",
			SportNFL,
			TrendingTypeAdd,
			TrendingPlayerOptions{LookbackHours: -1},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			players, err := testClient.ListTrendingPlayers(context.Background(), tc.sport, tc.trendingType, tc.options)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got players: %v", players)
				return
			}

			if len(players) == 0 {
				t.Errorf("expected non-empty player list")
				return
			}

			t.Logf("retrieved %d trending players", len(players))
		})
	}
}
