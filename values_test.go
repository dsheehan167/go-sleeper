package sleeper

import (
	"context"
	"testing"
)

func TestPlayerValues_Get(t *testing.T) {
	tt := []struct {
		testcase   string
		sport      Sport
		season     string
		format     ScoringFormat
		options    PlayerValuesOptions
		shouldPass bool
	}{
		{
			"valid - nfl half ppr",
			SportNFL,
			"2026",
			ScoringFormatHalfPPR,
			PlayerValuesOptions{},
			true,
		},
		{
			"valid - nfl dynasty ppr",
			SportNFL,
			"2026",
			ScoringFormatDynastyPPR,
			PlayerValuesOptions{IsDynasty: true},
			true,
		},
		{
			"valid - nfl standard with idp players",
			SportNFL,
			"2026",
			ScoringFormatStandard,
			PlayerValuesOptions{IncludeIDP: true},
			true,
		},
		{
			"missing sport",
			"",
			"2026",
			ScoringFormatHalfPPR,
			PlayerValuesOptions{},
			false,
		},
		{
			"missing season",
			SportNFL,
			"",
			ScoringFormatHalfPPR,
			PlayerValuesOptions{},
			false,
		},
		{
			"missing format",
			SportNFL,
			"2026",
			"",
			PlayerValuesOptions{},
			false,
		},
		{
			"unknown format returns empty result",
			SportNFL,
			"2026",
			"not_a_format",
			PlayerValuesOptions{},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			values, err := testClient.GetPlayerValues(context.Background(), tc.sport, tc.season, tc.format, tc.options)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d player values", len(values))
				return
			}

			t.Logf("retrieved %d player values", len(values))
		})
	}
}
