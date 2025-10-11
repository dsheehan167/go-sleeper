package sleeper

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	playerNicknamePrefix = "p_nick_"
)

// Roster represents a fantasy team roster in a Sleeper league.
type Roster struct {
	CoOwners []string        `json:"co_owners,omitempty"`
	Keepers  []string        `json:"keepers,omitempty"`
	LeagueID string          `json:"league_id,omitempty"`
	Metadata *RosterMetadata `json:"metadata,omitempty"`
	OwnerID  string          `json:"owner_id,omitempty"`
	Players  []string        `json:"players,omitempty"`
	Reserve  []string        `json:"reserve,omitempty"`
	RosterID int             `json:"roster_id,omitempty"`
	Settings *RosterSettings `json:"settings,omitempty"`
	Starters []string        `json:"starters,omitempty"`
	Taxi     []string        `json:"taxi,omitempty"`
}

// RosterSettings contains scoring and record settings for a roster.
type RosterSettings struct {
	Fpts               int `json:"fpts,omitempty"`
	FptsAgainst        int `json:"fpts_against,omitempty"`
	FptsAgainstDecimal int `json:"fpts_against_decimal,omitempty"`
	FptsDecimal        int `json:"fpts_decimal,omitempty"`
	Losses             int `json:"losses,omitempty"`
	Ppts               int `json:"ppts,omitempty"`
	PptsDecimal        int `json:"ppts_decimal,omitempty"`
	Ties               int `json:"ties,omitempty"`
	TotalMoves         int `json:"total_moves,omitempty"`
	WaiverBudgetUsed   int `json:"waiver_budget_used,omitempty"`
	WaiverPosition     int `json:"waiver_position,omitempty"`
	Wins               int `json:"wins,omitempty"`
}

// RosterMetadata contains metadata and player nicknames for a roster.
type RosterMetadata struct {
	AllowPnNews                   string            `json:"allow_pn_news,omitempty"`
	AllowPnScoring                string            `json:"allow_pn_scoring,omitempty"`
	AllowPnInactiveStarters       string            `json:"allow_pn_inactive_starters,omitempty"`
	AllowPnPlayerInjuryStatus     string            `json:"allow_pn_player_injury_status,omitempty"`
	RestrictPnScoringStartersOnly string            `json:"restrict_pn_scoring_starters_only,omitempty"`
	Record                        string            `json:"record,omitempty"`
	Streak                        string            `json:"streak,omitempty"`
	PlayerNicknames               map[string]string `json:"-"`
}

// UnmarshalJSON handles both regular fields and dynamic playerNickname fields
func (rm *RosterMetadata) UnmarshalJSON(data []byte) error {
	type Alias RosterMetadata
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(rm),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("unmarshalling roster metadata: %w", err)
	}

	// Use map[string]interface{} to handle any value type
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshalling raw roster metadata: %w", err)
	}

	// Initialize PlayerNicknames map and populate from p_nick_ fields
	rm.PlayerNicknames = make(map[string]string)
	for key, value := range raw {
		if strings.HasPrefix(key, playerNicknamePrefix) {
			playerID := strings.TrimPrefix(key, playerNicknamePrefix)
			// Only add if the value is actually a string
			if nickname, ok := value.(string); ok {
				rm.PlayerNicknames[playerID] = nickname
			}
		}
	}

	return nil
}
