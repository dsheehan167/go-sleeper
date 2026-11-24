package sleeper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestUnitProjections_List(t *testing.T) {
	tt := []struct {
		testcase      string
		sport         Sport
		season        string
		options       ProjectionOptions
		wantPath      string
		wantQuery     string
		status        int
		body          string
		expectedCount int
		expectedADP   float64
		wantErr       string
		shouldPass    bool
	}{
		{
			"defaults to regular season",
			SportNFL,
			"2026",
			ProjectionOptions{},
			"/projections/nfl/2026",
			"season_type=regular",
			http.StatusOK,
			`[{"player_id":"4034","stats":{"adp_ppr":1.4,"pts_ppr":350.2}},{"player_id":"9999","stats":{"adp_ppr":999}}]`,
			2,
			1.4,
			"",
			true,
		},
		{
			"all options",
			SportNFL,
			"2026",
			ProjectionOptions{Positions: []string{"QB", "RB"}, SeasonType: "post", OrderBy: "adp_ppr"},
			"/projections/nfl/2026",
			"order_by=adp_ppr&position%5B%5D=QB&position%5B%5D=RB&season_type=post",
			http.StatusOK,
			`[{"player_id":"4034","stats":{"adp_ppr":2},"player":{"player_id":"4034","position":"RB"}}]`,
			1,
			2,
			"",
			true,
		},
		{
			"missing fields",
			"",
			"",
			ProjectionOptions{},
			"",
			"",
			0,
			"",
			0,
			0,
			"sport is required\nseason is required",
			false,
		},
		{
			"server error",
			SportNFL,
			"2026",
			ProjectionOptions{},
			"/projections/nfl/2026",
			"season_type=regular",
			http.StatusInternalServerError,
			``,
			0,
			0,
			"getting projections: API error: 500",
			false,
		},
		{
			"bad json",
			SportNFL,
			"2026",
			ProjectionOptions{},
			"/projections/nfl/2026",
			"season_type=regular",
			http.StatusOK,
			`[{"stats":{"adp_ppr":"first"}}]`,
			0,
			0,
			"unmarshaling projections",
			false,
		},
		{
			"empty list",
			Sport("xfl"),
			"2026",
			ProjectionOptions{},
			"/projections/xfl/2026",
			"season_type=regular",
			http.StatusOK,
			`[]`,
			0,
			0,
			"projections not found",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respondExact(t, "api.sleeper.com", tc.wantPath, tc.wantQuery, tc.status, tc.body))

			projections, err := client.ListProjections(context.Background(), tc.sport, tc.season, tc.options)
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
				t.Errorf("expected failure but got %d projections", len(projections))
				return
			}

			if len(projections) != tc.expectedCount {
				t.Errorf("expected %d projections, got %d", tc.expectedCount, len(projections))
				return
			}

			if adp := projections[0].Stats.ADP(ADPTypePPR); adp != tc.expectedADP {
				t.Errorf("expected first ppr adp %v, got %v", tc.expectedADP, adp)
				return
			}
		})
	}
}

func TestUnitProjectionStats_ADP(t *testing.T) {
	stats := ProjectionStats{
		ADPStandard:        1,
		ADPHalfPPR:         2,
		ADPPPR:             3,
		ADP2QB:             4,
		ADPDynasty:         5,
		ADPDynastyStandard: 6,
		ADPDynastyHalfPPR:  7,
		ADPDynastyPPR:      8,
		ADPDynasty2QB:      9,
		ADPIDP:             10,
		ADPIDP1QB:          11,
		ADPRookie:          12,
	}

	tt := []struct {
		testcase string
		adpType  ADPType
		expected float64
	}{
		{"standard", ADPTypeStandard, 1},
		{"half ppr", ADPTypeHalfPPR, 2},
		{"ppr", ADPTypePPR, 3},
		{"2qb", ADPType2QB, 4},
		{"dynasty", ADPTypeDynasty, 5},
		{"dynasty standard", ADPTypeDynastyStandard, 6},
		{"dynasty half ppr", ADPTypeDynastyHalfPPR, 7},
		{"dynasty ppr", ADPTypeDynastyPPR, 8},
		{"dynasty 2qb", ADPTypeDynasty2QB, 9},
		{"idp", ADPTypeIDP, 10},
		{"idp 1qb", ADPTypeIDP1QB, 11},
		{"rookie has no adp type", ADPType("adp_rookie"), 0},
		{"unknown type", ADPType("adp_unknown"), 0},
		{"empty type", ADPType(""), 0},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			if adp := stats.ADP(tc.adpType); adp != tc.expected {
				t.Errorf("expected adp %v, got %v", tc.expected, adp)
				return
			}
		})
	}
}
