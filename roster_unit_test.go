package sleeper

import (
	"encoding/json"
	"maps"
	"testing"
)

func TestUnitRosterMetadata_Unmarshal(t *testing.T) {
	tt := []struct {
		testcase          string
		data              string
		expectedRecord    string
		expectedStreak    string
		expectedNicknames map[string]string
		shouldPass        bool
	}{
		{
			"regular fields and nicknames",
			`{"record":"WLW","streak":"1W","allow_pn_news":"on","p_nick_4034":"CMC","p_nick_6794":"JJ"}`,
			"WLW",
			"1W",
			map[string]string{"4034": "CMC", "6794": "JJ"},
			true,
		},
		{
			"no nicknames",
			`{"record":"LL"}`,
			"LL",
			"",
			map[string]string{},
			true,
		},
		{
			"empty object",
			`{}`,
			"",
			"",
			map[string]string{},
			true,
		},
		{
			"non string nicknames are ignored",
			`{"p_nick_1":"Ace","p_nick_2":7,"p_nick_3":null,"p_nick_4":{"a":1}}`,
			"",
			"",
			map[string]string{"1": "Ace"},
			true,
		},
		{
			"empty nickname is kept",
			`{"p_nick_1":""}`,
			"",
			"",
			map[string]string{"1": ""},
			true,
		},
		{
			"unknown keys are ignored",
			`{"nickname":"x","p_nick":"y","record":"W"}`,
			"W",
			"",
			map[string]string{},
			true,
		},
		{
			"wrong type for regular field",
			`{"record":5}`,
			"",
			"",
			nil,
			false,
		},
		{
			"invalid json",
			`{"record":`,
			"",
			"",
			nil,
			false,
		},
		{
			"array instead of object",
			`["record"]`,
			"",
			"",
			nil,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			var metadata RosterMetadata
			err := metadata.UnmarshalJSON([]byte(tc.data))
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %+v", metadata)
				return
			}

			if metadata.Record != tc.expectedRecord {
				t.Errorf("expected record %q, got %q", tc.expectedRecord, metadata.Record)
			}
			if metadata.Streak != tc.expectedStreak {
				t.Errorf("expected streak %q, got %q", tc.expectedStreak, metadata.Streak)
			}
			if !maps.Equal(metadata.PlayerNicknames, tc.expectedNicknames) {
				t.Errorf("expected nicknames %v, got %v", tc.expectedNicknames, metadata.PlayerNicknames)
			}
		})
	}
}

func TestUnitRoster_Unmarshal(t *testing.T) {
	tt := []struct {
		testcase          string
		data              string
		expectedRosterID  int
		expectedNicknames map[string]string
		shouldPass        bool
	}{
		{
			"roster with metadata",
			`{"roster_id":3,"owner_id":"99","players":["4034"],"metadata":{"record":"W","p_nick_4034":"CMC"}}`,
			3,
			map[string]string{"4034": "CMC"},
			true,
		},
		{
			"roster without metadata",
			`{"roster_id":4}`,
			4,
			nil,
			true,
		},
		{
			"roster with invalid metadata",
			`{"roster_id":5,"metadata":{"streak":[]}}`,
			0,
			nil,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			var roster Roster
			err := json.Unmarshal([]byte(tc.data), &roster)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %+v", roster)
				return
			}

			if roster.RosterID != tc.expectedRosterID {
				t.Errorf("expected roster id %d, got %d", tc.expectedRosterID, roster.RosterID)
			}

			if tc.expectedNicknames == nil {
				if roster.Metadata != nil {
					t.Errorf("expected nil metadata, got %+v", roster.Metadata)
				}
				return
			}

			if roster.Metadata == nil {
				t.Errorf("expected metadata, got nil")
				return
			}
			if !maps.Equal(roster.Metadata.PlayerNicknames, tc.expectedNicknames) {
				t.Errorf("expected nicknames %v, got %v", tc.expectedNicknames, roster.Metadata.PlayerNicknames)
			}
		})
	}
}
