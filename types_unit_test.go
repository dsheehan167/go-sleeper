package sleeper

import (
	"encoding/json"
	"testing"
)

func TestUnitFlexibleString_Unmarshal(t *testing.T) {
	tt := []struct {
		testcase   string
		data       string
		expected   FlexibleString
		shouldPass bool
	}{
		{
			"string",
			`"2024"`,
			"2024",
			true,
		},
		{
			"empty string",
			`""`,
			"",
			true,
		},
		{
			"int",
			`2024`,
			"2024",
			true,
		},
		{
			"negative int",
			`-7`,
			"-7",
			true,
		},
		{
			"null leaves zero value",
			`null`,
			"",
			true,
		},
		{
			"float",
			`2024.5`,
			"",
			false,
		},
		{
			"bool",
			`true`,
			"",
			false,
		},
		{
			"object",
			`{"season":"2024"}`,
			"",
			false,
		},
		{
			"array",
			`["2024"]`,
			"",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			var f FlexibleString
			err := f.UnmarshalJSON([]byte(tc.data))
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %q", f)
				return
			}

			if f != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, f)
				return
			}

			if f.String() != string(tc.expected) {
				t.Errorf("expected String() %q, got %q", tc.expected, f.String())
				return
			}
		})
	}
}

func TestUnitFlexibleString_InStruct(t *testing.T) {
	tt := []struct {
		testcase   string
		data       string
		expected   FlexibleString
		shouldPass bool
	}{
		{
			"previous season as string",
			`{"previous_season":"2023"}`,
			"2023",
			true,
		},
		{
			"previous season as int",
			`{"previous_season":2023}`,
			"2023",
			true,
		},
		{
			"previous season missing",
			`{}`,
			"",
			true,
		},
		{
			"previous season invalid",
			`{"previous_season":false}`,
			"",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			var state SportState
			err := json.Unmarshal([]byte(tc.data), &state)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %+v", state)
				return
			}

			if state.PreviousSeason != tc.expected {
				t.Errorf("expected previous season %q, got %q", tc.expected, state.PreviousSeason)
				return
			}
		})
	}
}

func TestUnitFlexibleString_Int(t *testing.T) {
	tt := []struct {
		testcase   string
		value      FlexibleString
		expected   int
		shouldPass bool
	}{
		{
			"numeric",
			"2024",
			2024,
			true,
		},
		{
			"negative",
			"-3",
			-3,
			true,
		},
		{
			"non numeric",
			"abc",
			0,
			false,
		},
		{
			"empty",
			"",
			0,
			false,
		},
		{
			"decimal",
			"1.5",
			0,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			i, err := tc.value.Int()
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if i != 0 {
					t.Errorf("expected 0 on error, got %d", i)
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d", i)
				return
			}

			if i != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, i)
				return
			}
		})
	}
}
