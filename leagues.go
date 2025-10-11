package sleeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// League represents a Sleeper fantasy league.
type League struct {
	Avatar                  string           `json:"avatar,omitempty"`
	Sport                   string           `json:"sport,omitempty"`
	LastMessageID           string           `json:"last_message_id,omitempty"`
	LastMessageTime         int64            `json:"last_message_time,omitempty"`
	Shard                   int              `json:"shard,omitempty"`
	LastTransactionID       int64            `json:"last_transaction_id,omitempty"`
	LastPinnedMessageID     string           `json:"last_pinned_message_id,omitempty"`
	PreviousLeagueID        string           `json:"previous_league_id,omitempty"`
	LoserBracketOverridesID int64            `json:"loser_bracket_overrides_id,omitempty"`
	RosterPositions         []string         `json:"roster_positions,omitempty"`
	LastAuthorIsBot         bool             `json:"last_author_is_bot,omitempty"`
	BracketID               int64            `json:"bracket_id,omitempty"`
	CompanyID               string           `json:"company_id,omitempty"`
	Status                  string           `json:"status,omitempty"`
	LoserBracketID          int64            `json:"loser_bracket_id,omitempty"`
	ScoringSettings         *ScoringSettings `json:"scoring_settings,omitempty"`
	DisplayOrder            int              `json:"display_order,omitempty"`
	BracketOverridesID      string           `json:"bracket_overrides_id,omitempty"`
	LastReadID              string           `json:"last_read_id,omitempty"`
	Name                    string           `json:"name,omitempty"`
	Season                  string           `json:"season,omitempty"`
	LastAuthorAvatar        string           `json:"last_author_avatar,omitempty"`
	LeagueID                string           `json:"league_id,omitempty"`
	DraftID                 string           `json:"draft_id,omitempty"`
	GroupID                 string           `json:"group_id,omitempty"`
	Metadata                *LeagueMetadata  `json:"metadata,omitempty"`
	LastMessageAttachment   string           `json:"last_message_attachment,omitempty"`
	LastAuthorDisplayName   string           `json:"last_author_display_name,omitempty"`
	LastAuthorID            string           `json:"last_author_id,omitempty"`
	TotalRosters            int              `json:"total_rosters,omitempty"`
	Settings                *Settings        `json:"settings,omitempty"`
	SeasonType              string           `json:"season_type,omitempty"`
	LastMessageTextMap      string           `json:"last_message_text_map,omitempty"`
}

// ScoringSettings holds the scoring rules for a league.
type ScoringSettings struct {
	Sack         float64 `json:"sack,omitempty"`
	FGM4049      float64 `json:"fgm_40_49,omitempty"`
	FGMYds       float64 `json:"fgm_yds,omitempty"`
	PassInt      float64 `json:"pass_int,omitempty"`
	PtsAllow0    float64 `json:"pts_allow_0,omitempty"`
	Pass2pt      float64 `json:"pass_2pt,omitempty"`
	StTd         float64 `json:"st_td,omitempty"`
	FGMYdsOver30 float64 `json:"fgm_yds_over_30,omitempty"`
	RecTd        float64 `json:"rec_td,omitempty"`
	FGM3039      float64 `json:"fgm_30_39,omitempty"`
	FGM5059      float64 `json:"fgm_50_59,omitempty"`
	XPMiss       float64 `json:"xpmiss,omitempty"`
	RushTd       float64 `json:"rush_td,omitempty"`
	DefPrTd      float64 `json:"def_pr_td,omitempty"`
	Def4AndStop  float64 `json:"def_4_and_stop,omitempty"`
	Rec2pt       float64 `json:"rec_2pt,omitempty"`
	PassIntTd    float64 `json:"pass_int_td,omitempty"`
	StFumRec     float64 `json:"st_fum_rec,omitempty"`
	FGMiss       float64 `json:"fgmiss,omitempty"`
	FF           float64 `json:"ff,omitempty"`
	Rec          float64 `json:"rec,omitempty"`
	PtsAllow1420 float64 `json:"pts_allow_14_20,omitempty"`
	FGM019       float64 `json:"fgm_0_19,omitempty"`
	DefKrTd      float64 `json:"def_kr_td,omitempty"`
	Int          float64 `json:"int,omitempty"`
	DefStFumRec  float64 `json:"def_st_fum_rec,omitempty"`
	FumLost      float64 `json:"fum_lost,omitempty"`
	PtsAllow16   float64 `json:"pts_allow_1_6,omitempty"`
	FGM2029      float64 `json:"fgm_20_29,omitempty"`
	PtsAllow2127 float64 `json:"pts_allow_21_27,omitempty"`
	XPM          float64 `json:"xpm,omitempty"`
	Rush2pt      float64 `json:"rush_2pt,omitempty"`
	FumRec       float64 `json:"fum_rec,omitempty"`
	DefStTd      float64 `json:"def_st_td,omitempty"`
	FGM50p       float64 `json:"fgm_50p,omitempty"`
	DefTd        float64 `json:"def_td,omitempty"`
	Safe         float64 `json:"safe,omitempty"`
	PassYd       float64 `json:"pass_yd,omitempty"`
	BlkKick      float64 `json:"blk_kick,omitempty"`
	PassTd       float64 `json:"pass_td,omitempty"`
	RushYd       float64 `json:"rush_yd,omitempty"`
	Fum          float64 `json:"fum,omitempty"`
	PtsAllow2834 float64 `json:"pts_allow_28_34,omitempty"`
	PtsAllow35p  float64 `json:"pts_allow_35p,omitempty"`
	FumRecTd     float64 `json:"fum_rec_td,omitempty"`
	RecYd        float64 `json:"rec_yd,omitempty"`
	DefStFF      float64 `json:"def_st_ff,omitempty"`
	PtsAllow713  float64 `json:"pts_allow_7_13,omitempty"`
	StFF         float64 `json:"st_ff,omitempty"`
}

// LeagueMetadata contains additional information about a league.
type LeagueMetadata struct {
	AutoContinue               string `json:"auto_continue,omitempty"`
	KeeperDeadline             string `json:"keeper_deadline,omitempty"`
	LatestLeagueWinnerRosterID string `json:"latest_league_winner_roster_id,omitempty"`
}

// Settings contains configuration options for a league.
type Settings struct {
	BestBall                 int `json:"best_ball,omitempty"`
	LastReport               int `json:"last_report,omitempty"`
	WaiverBudget             int `json:"waiver_budget,omitempty"`
	DisableAdds              int `json:"disable_adds,omitempty"`
	CapacityOverride         int `json:"capacity_override,omitempty"`
	TaxiDeadline             int `json:"taxi_deadline,omitempty"`
	DraftRounds              int `json:"draft_rounds,omitempty"`
	ReserveAllowNA           int `json:"reserve_allow_na,omitempty"`
	StartWeek                int `json:"start_week,omitempty"`
	PlayoffSeedType          int `json:"playoff_seed_type,omitempty"`
	PlayoffTeams             int `json:"playoff_teams,omitempty"`
	VetoVotesNeeded          int `json:"veto_votes_needed,omitempty"`
	Squads                   int `json:"squads,omitempty"`
	NumTeams                 int `json:"num_teams,omitempty"`
	DailyWaiversHour         int `json:"daily_waivers_hour,omitempty"`
	PlayoffType              int `json:"playoff_type,omitempty"`
	TaxiSlots                int `json:"taxi_slots,omitempty"`
	SubStartTimeEligibility  int `json:"sub_start_time_eligibility,omitempty"`
	LastScoredLeg            int `json:"last_scored_leg,omitempty"`
	DailyWaiversDays         int `json:"daily_waivers_days,omitempty"`
	SubLockIfStarterActive   int `json:"sub_lock_if_starter_active,omitempty"`
	PlayoffWeekStart         int `json:"playoff_week_start,omitempty"`
	WaiverClearDays          int `json:"waiver_clear_days,omitempty"`
	ReserveAllowDoubtful     int `json:"reserve_allow_doubtful,omitempty"`
	CommissionerDirectInvite int `json:"commissioner_direct_invite,omitempty"`
	VetoAutoPoll             int `json:"veto_auto_poll,omitempty"`
	ReserveAllowDNR          int `json:"reserve_allow_dnr,omitempty"`
	TaxiAllowVets            int `json:"taxi_allow_vets,omitempty"`
	WaiverDayOfWeek          int `json:"waiver_day_of_week,omitempty"`
	PlayoffRoundType         int `json:"playoff_round_type,omitempty"`
	ReserveAllowOut          int `json:"reserve_allow_out,omitempty"`
	ReserveAllowSus          int `json:"reserve_allow_sus,omitempty"`
	VetoShowVotes            int `json:"veto_show_votes,omitempty"`
	TradeDeadline            int `json:"trade_deadline,omitempty"`
	TaxiYears                int `json:"taxi_years,omitempty"`
	DailyWaivers             int `json:"daily_waivers,omitempty"`
	FaabSuggestions          int `json:"faab_suggestions,omitempty"`
	DisableTrades            int `json:"disable_trades,omitempty"`
	PickTrading              int `json:"pick_trading,omitempty"`
	Type                     int `json:"type,omitempty"`
	MaxKeepers               int `json:"max_keepers,omitempty"`
	WaiverType               int `json:"waiver_type,omitempty"`
	MaxSubs                  int `json:"max_subs,omitempty"`
	LeagueAverageMatch       int `json:"league_average_match,omitempty"`
	TradeReviewDays          int `json:"trade_review_days,omitempty"`
	BenchLock                int `json:"bench_lock,omitempty"`
	OffseasonAdds            int `json:"offseason_adds,omitempty"`
	Leg                      int `json:"leg,omitempty"`
	ReserveSlots             int `json:"reserve_slots,omitempty"`
	ReserveAllowCov          int `json:"reserve_allow_cov,omitempty"`
	DailyWaiversLastRan      int `json:"daily_waivers_last_ran,omitempty"`
}

// Matchup represents a weekly matchup in a league.
type Matchup struct {
	Points         float64            `json:"points,omitempty"`
	Players        []string           `json:"players,omitempty"`
	RosterID       int                `json:"roster_id,omitempty"`
	CustomPoints   *float64           `json:"custom_points,omitempty"`
	MatchupID      int                `json:"matchup_id,omitempty"`
	Starters       []string           `json:"starters,omitempty"`
	StartersPoints []float64          `json:"starters_points,omitempty"`
	PlayersPoints  map[string]float64 `json:"players_points,omitempty"`
}

// PlayoffMatchup represents a playoff bracket matchup.
type PlayoffMatchup struct {
	Round     int                 `json:"r,omitempty"`       // The round for this matchup, 1st, 2nd, 3rd round, etc.
	MatchID   int                 `json:"m,omitempty"`       // The match id of the matchup, unique for all matchups within a bracket.
	Team1     *int                `json:"t1,omitempty"`      // The roster_id of a team in this matchup OR {w: 1} which means the winner of match id 1.
	Team2     *int                `json:"t2,omitempty"`      // The roster_id of the other team in this matchup OR {l: 1} which means the loser of match id 1.
	Winner    *int                `json:"w,omitempty"`       // The roster_id of the winning team, if the match has been played.
	Loser     *int                `json:"l,omitempty"`       // The roster_id of the losing team, if the match has been played.
	Team1From *PlayoffMatchupFrom `json:"t1_from,omitempty"` // Where t1 comes from, either winner or loser of the match id, necessary to show bracket progression.
	Team2From *PlayoffMatchupFrom `json:"t2_from,omitempty"` // Where t2 comes from, either winner or loser of the match id, necessary to show bracket progression.
	Placement *int                `json:"p,omitempty"`       // Position/placement (e.g., 1st place, 3rd place, 5th place).
}

// PlayoffMatchupFrom describes where a playoff team comes from in the bracket.
type PlayoffMatchupFrom struct {
	WinnerOfMatch *int `json:"w,omitempty"` // Winner of match id
	LoserOfMatch  *int `json:"l,omitempty"` // Loser of match id
}

// GetLeague retrieves a League by league ID.
func (c *Client) GetLeague(ctx context.Context, leagueID string) (*League, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, errors.New("leagueID is required")
	}

	endpoint := c.buildEndpoint(endpointLeague, leagueID)

	var league *League
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting league: %w", err)
	}

	if err := json.Unmarshal(by, &league); err != nil {
		return nil, fmt.Errorf("unmarshaling league: %w", err)
	}

	return league, nil
}

// GetLeagueRosters retrieves all rosters for a given league ID.
func (c *Client) GetLeagueRosters(ctx context.Context, leagueID string) ([]*Roster, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, errors.New("leagueID is required")
	}

	endpoint := c.buildEndpoint(endpointLeagueRosters, leagueID)

	var rosters []*Roster
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting league rosters: %w", err)
	}

	if err := json.Unmarshal(by, &rosters); err != nil {
		return nil, fmt.Errorf("unmarshaling league rosters: %w", err)
	}

	return rosters, nil
}

// GetLeagueUsers retrieves all users for a given league ID.
func (c *Client) GetLeagueUsers(ctx context.Context, leagueID string) ([]*User, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, errors.New("leagueID is required")
	}

	endpoint := c.buildEndpoint(endpointLeagueUsers, leagueID)

	var users []*User
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting league users: %w", err)
	}

	if err := json.Unmarshal(by, &users); err != nil {
		return nil, fmt.Errorf("unmarshaling league users: %w", err)
	}

	return users, nil
}

// GetLeagueMatchups retrieves all matchups for a given league ID and week.
func (c *Client) GetLeagueMatchups(ctx context.Context, leagueID string, week int) ([]*Matchup, error) {

	var errs []string
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		errs = append(errs, "leagueID is required")
	}
	if week < 1 {
		errs = append(errs, "week must be greater than or equal to 1")
	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	endpoint := c.buildEndpoint(endpointLeagueMatchups, leagueID, week)

	var matchups []*Matchup
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting league matchups: %w", err)
	}

	if err := json.Unmarshal(by, &matchups); err != nil {
		return nil, fmt.Errorf("unmarshaling league matchups: %w", err)
	}

	return matchups, nil
}

// GetTransactions retrieves all transactions for a given league ID and week.
func (c *Client) GetTransactions(ctx context.Context, leagueID string, week int) ([]*Transaction, error) {
	var errs []string
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		errs = append(errs, "leagueID is required")
	}
	if week < 1 {
		errs = append(errs, "week must be greater than or equal to 1")
	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	endpoint := c.buildEndpoint(endpointLeagueTransactions, leagueID, week)

	var transactions []*Transaction
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting league transactions: %w", err)
	}

	if err := json.Unmarshal(by, &transactions); err != nil {
		return nil, fmt.Errorf("unmarshaling league transactions: %w", err)
	}

	return transactions, nil
}

// GetLeagueTradedPicks retrieves all traded draft picks for a given league ID.
func (c *Client) GetLeagueTradedPicks(ctx context.Context, leagueID string) ([]*TradedDraftPick, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, errors.New("leagueID is required")
	}

	endpoint := c.buildEndpoint(endpointLeagueTradedPicks, leagueID)

	var picks []*TradedDraftPick
	data, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting traded picks: %w", err)
	}

	if err := json.Unmarshal(data, &picks); err != nil {
		return nil, fmt.Errorf("unmarshaling traded picks: %w", err)
	}

	return picks, nil
}

// GetLeagueWinnersBracket retrieves the winners bracket for a given league ID.
func (c *Client) GetLeagueWinnersBracket(ctx context.Context, leagueID string) ([]*PlayoffMatchup, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, errors.New("leagueID is required")
	}

	endpoint := c.buildEndpoint(endpointLeagueWinnersBracket, leagueID)

	var matchups []*PlayoffMatchup
	data, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting winners bracket: %w", err)
	}

	if err := json.Unmarshal(data, &matchups); err != nil {
		return nil, fmt.Errorf("unmarshaling winners bracket: %w", err)
	}

	return matchups, nil
}

// GetLeagueLosersBracket retrieves the losers bracket for a given league ID.
func (c *Client) GetLeagueLosersBracket(ctx context.Context, leagueID string) ([]*PlayoffMatchup, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, errors.New("leagueID is required")
	}

	endpoint := c.buildEndpoint(endpointLeagueLosersBracket, leagueID)

	var matchups []*PlayoffMatchup
	data, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting losers bracket: %w", err)
	}

	if err := json.Unmarshal(data, &matchups); err != nil {
		return nil, fmt.Errorf("unmarshaling losers bracket: %w", err)
	}

	return matchups, nil
}
