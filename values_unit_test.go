package sleeper

import (
	"net/http"
	"strings"
	"testing"
)

func TestUnitPlayerValues_Get(t *testing.T) {
	tt := []struct {
		testcase      string
		sport         Sport
		season        string
		format        ScoringFormat
		options       PlayerValuesOptions
		wantPath      string
		wantQuery     string
		status        int
		body          string
		expectedCount int
		wantErr       string
		shouldPass    bool
	}{
		{
			"valid nfl half ppr",
			SportNFL,
			"2026",
			ScoringFormatHalfPPR,
			PlayerValuesOptions{},
			"/players/nfl/values/regular/2026/half_ppr",
			"idp=false&is_dynasty=false",
			http.StatusOK,
			`{"4034":9850.5,"6794":9100}`,
			2,
			"",
			true,
		},
		{
			"valid with options",
			SportNFL,
			"2026",
			ScoringFormatDynastyPPR,
			PlayerValuesOptions{IncludeIDP: true, IsDynasty: true},
			"/players/nfl/values/regular/2026/dynasty_ppr",
			"idp=true&is_dynasty=true",
			http.StatusOK,
			`{"4034":9850.5}`,
			1,
			"",
			true,
		},
		{
			"missing fields",
			"",
			"",
			"",
			PlayerValuesOptions{},
			"",
			"",
			0,
			"",
			0,
			"sport is required\nseason is required\nformat is required",
			false,
		},
		{
			"missing format",
			SportNFL,
			"2026",
			"",
			PlayerValuesOptions{},
			"",
			"",
			0,
			"",
			0,
			"format is required",
			false,
		},
		{
			"server error",
			SportNFL,
			"2026",
			ScoringFormatPPR,
			PlayerValuesOptions{},
			"/players/nfl/values/regular/2026/ppr",
			"idp=false&is_dynasty=false",
			http.StatusInternalServerError,
			``,
			0,
			"getting player values: API error: 500",
			false,
		},
		{
			"bad json",
			SportNFL,
			"2026",
			ScoringFormatPPR,
			PlayerValuesOptions{},
			"/players/nfl/values/regular/2026/ppr",
			"idp=false&is_dynasty=false",
			http.StatusOK,
			`{"4034":"high"}`,
			0,
			"unmarshaling player values",
			false,
		},
		{
			"empty object",
			SportNFL,
			"1999",
			ScoringFormatPPR,
			PlayerValuesOptions{},
			"/players/nfl/values/regular/1999/ppr",
			"idp=false&is_dynasty=false",
			http.StatusOK,
			`{}`,
			0,
			"player values not found",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respondExact(t, "api.sleeper.com", tc.wantPath, tc.wantQuery, tc.status, tc.body))

			values, err := client.GetPlayerValues(t.Context(), tc.sport, tc.season, tc.format, tc.options)
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
				t.Errorf("expected failure but got %d values", len(values))
				return
			}

			if len(values) != tc.expectedCount {
				t.Errorf("expected %d values, got %d", tc.expectedCount, len(values))
				return
			}
		})
	}
}
