package sleeper

import (
	"context"
	"testing"
)

func TestProjections_List(t *testing.T) {
	tt := []struct {
		testcase   string
		sport      Sport
		season     string
		options    ProjectionOptions
		shouldPass bool
	}{
		{
			"valid - nfl all positions",
			SportNFL,
			"2026",
			ProjectionOptions{},
			true,
		},
		{
			"valid - nfl filtered positions ordered by adp",
			SportNFL,
			"2026",
			ProjectionOptions{
				Positions: []string{"QB", "RB"},
				OrderBy:   string(ADPTypeDynasty2QB),
			},
			true,
		},
		{
			"missing sport",
			"",
			"2026",
			ProjectionOptions{},
			false,
		},
		{
			"missing season",
			SportNFL,
			"",
			ProjectionOptions{},
			false,
		},
		{
			"unknown sport returns empty result",
			"foosball",
			"2026",
			ProjectionOptions{},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			projections, err := testClient.ListProjections(context.Background(), tc.sport, tc.season, tc.options)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d projections", len(projections))
				return
			}

			var drafted, withPts int
			for _, p := range projections {
				if p.Stats.ADPDynasty2QB > 0 && p.Stats.ADPDynasty2QB < ADPUndrafted {
					drafted++
				}
				if p.Stats.PtsPPR > 0 {
					withPts++
				}
			}
			if drafted == 0 || withPts == 0 {
				t.Errorf("stats did not decode: %d drafted, %d with projected points", drafted, withPts)
				return
			}

			t.Logf("retrieved %d projections (%d drafted, %d with projected points)", len(projections), drafted, withPts)
		})
	}
}
