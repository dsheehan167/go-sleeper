package sleeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Draft represents a draft object in the Sleeper API.
type Draft struct {
	Created         int64          `json:"created"`
	Creators        []string       `json:"creators"`
	DraftID         string         `json:"draft_id"`
	DraftOrder      map[string]int `json:"draft_order"`
	LastMessageID   string         `json:"last_message_id"`
	LastMessageTime int64          `json:"last_message_time"`
	LastPicked      *int64         `json:"last_picked"`
	LeagueID        string         `json:"league_id"`
	Metadata        *DraftMetadata `json:"metadata"`
	Season          string         `json:"season"`
	SeasonType      string         `json:"season_type"`
	Settings        *DraftSettings `json:"settings"`
	SlotToRosterID  map[string]int `json:"slot_to_roster_id"`
	Sport           string         `json:"sport"`
	StartTime       *int64         `json:"start_time"`
	Status          string         `json:"status"`
	Type            string         `json:"type"`
}

// DraftMetadata contains metadata for a draft in the Sleeper API.
type DraftMetadata struct {
	Description      string `json:"description"`
	Name             string `json:"name"`
	ScoringType      string `json:"scoring_type"`
	ShowTeamNames    string `json:"show_team_names,omitempty"`
	LeagueType       string `json:"league_type,omitempty"`
	ElapsedPickTimer string `json:"elapsed_pick_timer,omitempty"`
	IsAutopaused     string `json:"is_autopaused,omitempty"`
}

// DraftSettings contains settings for a draft in the Sleeper API.
type DraftSettings struct {
	AlphaSort             int `json:"alpha_sort"`
	AutopauseEnabled      int `json:"autopause_enabled"`
	AutopauseEndTime      int `json:"autopause_end_time"`
	AutopauseStartTime    int `json:"autopause_start_time"`
	Autostart             int `json:"autostart"`
	CPUAutopick           int `json:"cpu_autopick"`
	EnforcePositionLimits int `json:"enforce_position_limits"`
	NominationTimer       int `json:"nomination_timer"`
	PickTimer             int `json:"pick_timer"`
	PlayerType            int `json:"player_type"`
	ReversalRound         int `json:"reversal_round"`
	Rounds                int `json:"rounds"`
	Teams                 int `json:"teams"`

	// Roster slots - NFL
	SlotsBN        int `json:"slots_bn,omitempty"`
	SlotsFlex      int `json:"slots_flex,omitempty"`
	SlotsQB        int `json:"slots_qb,omitempty"`
	SlotsRB        int `json:"slots_rb,omitempty"`
	SlotsSuperFlex int `json:"slots_super_flex,omitempty"`
	SlotsTE        int `json:"slots_te,omitempty"`
	SlotsWR        int `json:"slots_wr,omitempty"`
	SlotsDEF       int `json:"slots_def,omitempty"`
	SlotsK         int `json:"slots_k,omitempty"`

	// Roster slots - NBA
	SlotsC    int `json:"slots_c,omitempty"`
	SlotsF    int `json:"slots_f,omitempty"`
	SlotsG    int `json:"slots_g,omitempty"`
	SlotsPF   int `json:"slots_pf,omitempty"`
	SlotsPG   int `json:"slots_pg,omitempty"`
	SlotsSF   int `json:"slots_sf,omitempty"`
	SlotsSG   int `json:"slots_sg,omitempty"`
	SlotsUtil int `json:"slots_util,omitempty"`
}

// DraftPick represents a single pick in a draft in the Sleeper API.
type DraftPick struct {
	DraftID   string              `json:"draft_id"`
	DraftSlot int                 `json:"draft_slot"`
	IsKeeper  interface{}         `json:"is_keeper"`
	Metadata  *DraftPickMetadata  `json:"metadata"`
	PickNo    int                 `json:"pick_no"`
	PickedBy  string              `json:"picked_by"`
	PlayerID  string              `json:"player_id"`
	Reactions map[string][]string `json:"reactions"`
	RosterID  int                 `json:"roster_id"`
	Round     int                 `json:"round"`
}

// DraftPickMetadata contains metadata for a draft pick in the Sleeper API.
type DraftPickMetadata struct {
	FirstName     string `json:"first_name"`
	InjuryStatus  string `json:"injury_status"`
	LastName      string `json:"last_name"`
	NewsUpdated   string `json:"news_updated"`
	Number        string `json:"number"`
	PlayerID      string `json:"player_id"`
	Position      string `json:"position"`
	Sport         string `json:"sport"`
	Status        string `json:"status"`
	Team          string `json:"team"`
	TeamAbbr      string `json:"team_abbr"`
	TeamChangedAt string `json:"team_changed_at"`
	YearsExp      string `json:"years_exp"`
}

// GetDraft retrieves a single draft by draft ID from the Sleeper API.
// See: https://docs.sleeper.com/#get-draft
func (c *Client) GetDraft(ctx context.Context, draftID string) (*Draft, error) {
	draftID = strings.TrimSpace(draftID)
	if draftID == "" {
		return nil, errors.New("draftID is required")
	}

	endpoint := c.buildEndpoint(endpointDraft, draftID)

	var draft *Draft
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting draft: %w", err)
	}

	if err := json.Unmarshal(by, &draft); err != nil {
		return nil, fmt.Errorf("unmarshaling draft: %w", err)
	}

	return draft, nil
}

// GetUserDrafts retrieves all drafts for a given user, sport, and season from the Sleeper API.
// See: https://docs.sleeper.com/#get-user-drafts
func (c *Client) GetUserDrafts(ctx context.Context, userID string, sport sport, season string) ([]*Draft, error) {
	userID = strings.TrimSpace(userID)
	season = strings.TrimSpace(season)
	var errs []string
	if userID == "" {
		errs = append(errs, "userID is required")
	}
	if sport == "" {
		errs = append(errs, "sport is required")
	}
	if season == "" {
		errs = append(errs, "season is required")
	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	endpoint := c.buildEndpoint(endpointUserDrafts, userID, sport, season)

	var drafts []*Draft
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting drafts for user: %w", err)
	}

	if err := json.Unmarshal(by, &drafts); err != nil {
		return nil, fmt.Errorf("unmarshaling drafts: %w", err)
	}

	return drafts, nil
}

// GetLeagueDrafts retrieves all drafts for a given league from the Sleeper API.
// See: https://docs.sleeper.com/#get-league-drafts
func (c *Client) GetLeagueDrafts(ctx context.Context, leagueID string) ([]*Draft, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, errors.New("leagueID is required")
	}

	endpoint := c.buildEndpoint(endpointLeagueDrafts, leagueID)

	var drafts []*Draft
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting drafts for league: %w", err)
	}

	if err := json.Unmarshal(by, &drafts); err != nil {
		return nil, fmt.Errorf("unmarshaling drafts: %w", err)
	}

	if len(drafts) == 0 {
		return nil, errors.New("no drafts found for league")
	}

	return drafts, nil
}

// GetDraftTradedPicks retrieves all traded picks for a given draft from the Sleeper API.
// See: https://docs.sleeper.com/#get-draft-traded-picks
func (c *Client) GetDraftTradedPicks(ctx context.Context, draftID string) ([]*TradedDraftPick, error) {
	draftID = strings.TrimSpace(draftID)
	if draftID == "" {
		return nil, errors.New("draftID is required")
	}

	endpoint := c.buildEndpoint(endpointDraftTradedPicks, draftID)

	var picks []*TradedDraftPick
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting traded picks for draft: %w", err)
	}

	if err := json.Unmarshal(by, &picks); err != nil {
		return nil, fmt.Errorf("unmarshaling traded picks: %w", err)
	}

	return picks, nil
}

// GetDraftPicks retrieves all picks for a given draft from the Sleeper API.
// See: https://docs.sleeper.com/#get-draft-picks
func (c *Client) GetDraftPicks(ctx context.Context, draftID string) ([]*DraftPick, error) {
	draftID = strings.TrimSpace(draftID)
	if draftID == "" {
		return nil, errors.New("draftID is required")
	}

	endpoint := c.buildEndpoint(endpointDraftPicks, draftID)

	var picks []*DraftPick
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting picks for draft: %w", err)
	}

	if err := json.Unmarshal(by, &picks); err != nil {
		return nil, fmt.Errorf("unmarshaling picks: %w", err)
	}

	return picks, nil
}
