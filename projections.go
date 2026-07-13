package sleeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ADPType selects a league format's average draft position (lower = earlier
// pick); use it with ProjectionStats.ADP or ProjectionOptions.OrderBy. These
// match the ADP shown in the Sleeper draft lobby, including for dynasty
// leagues. All types except ADPTypeDynasty are populated every season.
// ADPTypeDynasty is erratic — Sleeper filled it for 2022 and 2025 but not
// 2023, 2024, or 2026 — so prefer the scoring-specific ADPTypeDynasty* types.
type ADPType string

const (
	ADPTypeStandard        ADPType = "adp_std"
	ADPTypeHalfPPR         ADPType = "adp_half_ppr"
	ADPTypePPR             ADPType = "adp_ppr"
	ADPType2QB             ADPType = "adp_2qb"
	ADPTypeDynasty         ADPType = "adp_dynasty"
	ADPTypeDynastyStandard ADPType = "adp_dynasty_std"
	ADPTypeDynastyHalfPPR  ADPType = "adp_dynasty_half_ppr"
	ADPTypeDynastyPPR      ADPType = "adp_dynasty_ppr"
	ADPTypeDynasty2QB      ADPType = "adp_dynasty_2qb"
	ADPTypeIDP             ADPType = "adp_idp"
	ADPTypeIDP1QB          ADPType = "adp_idp_1qb"
)

// ADPUndrafted is the sentinel ADP reported for players not being drafted in a
// format; filter ADP values below it to get real draft positions.
const ADPUndrafted = 999.0

// ProjectionStats holds a player's projected season stats and per-format ADP.
// Only the fields relevant to the player's position are populated; absent
// stats decode to zero. Field set observed across the 2022-2026 seasons.
type ProjectionStats struct {
	// Average draft position per league format (999 = not being drafted).
	// ADPRookie exists in responses but has never been populated.
	ADPStandard        float64 `json:"adp_std,omitempty"`
	ADPHalfPPR         float64 `json:"adp_half_ppr,omitempty"`
	ADPPPR             float64 `json:"adp_ppr,omitempty"`
	ADP2QB             float64 `json:"adp_2qb,omitempty"`
	ADPDynasty         float64 `json:"adp_dynasty,omitempty"`
	ADPDynastyStandard float64 `json:"adp_dynasty_std,omitempty"`
	ADPDynastyHalfPPR  float64 `json:"adp_dynasty_half_ppr,omitempty"`
	ADPDynastyPPR      float64 `json:"adp_dynasty_ppr,omitempty"`
	ADPDynasty2QB      float64 `json:"adp_dynasty_2qb,omitempty"`
	ADPIDP             float64 `json:"adp_idp,omitempty"`
	ADPIDP1QB          float64 `json:"adp_idp_1qb,omitempty"`
	ADPRookie          float64 `json:"adp_rookie,omitempty"`

	// Projected fantasy points per scoring format
	GamesPlayed float64 `json:"gp,omitempty"`
	PtsStandard float64 `json:"pts_std,omitempty"`
	PtsHalfPPR  float64 `json:"pts_half_ppr,omitempty"`
	PtsPPR      float64 `json:"pts_ppr,omitempty"`

	// Passing
	PassAtt        float64 `json:"pass_att,omitempty"`
	PassCmp        float64 `json:"pass_cmp,omitempty"`
	PassCmpPct     float64 `json:"cmp_pct,omitempty"`
	PassYd         float64 `json:"pass_yd,omitempty"`
	PassTD         float64 `json:"pass_td,omitempty"`
	PassInt        float64 `json:"pass_int,omitempty"`
	PassIntTD      float64 `json:"pass_int_td,omitempty"`
	PassFirstDowns float64 `json:"pass_fd,omitempty"`
	Pass2Pt        float64 `json:"pass_2pt,omitempty"`

	// Rushing
	RushAtt        float64 `json:"rush_att,omitempty"`
	RushYd         float64 `json:"rush_yd,omitempty"`
	RushTD         float64 `json:"rush_td,omitempty"`
	RushFirstDowns float64 `json:"rush_fd,omitempty"`
	Rush2Pt        float64 `json:"rush_2pt,omitempty"`

	// Receiving (rec_N_M buckets are receptions by yardage gained)
	Rec           float64 `json:"rec,omitempty"`
	RecYd         float64 `json:"rec_yd,omitempty"`
	RecTD         float64 `json:"rec_td,omitempty"`
	RecFirstDowns float64 `json:"rec_fd,omitempty"`
	Rec2Pt        float64 `json:"rec_2pt,omitempty"`
	Rec0To4       float64 `json:"rec_0_4,omitempty"`
	Rec5To9       float64 `json:"rec_5_9,omitempty"`
	Rec10To19     float64 `json:"rec_10_19,omitempty"`
	Rec20To29     float64 `json:"rec_20_29,omitempty"`
	Rec30To39     float64 `json:"rec_30_39,omitempty"`
	Rec40Plus     float64 `json:"rec_40p,omitempty"`

	// Reception bonuses by position
	BonusRecRB float64 `json:"bonus_rec_rb,omitempty"`
	BonusRecTE float64 `json:"bonus_rec_te,omitempty"`
	BonusRecWR float64 `json:"bonus_rec_wr,omitempty"`

	// Fumbles
	FumLost float64 `json:"fum_lost,omitempty"`

	// Kicking
	XPM          float64 `json:"xpm,omitempty"`
	XPMiss       float64 `json:"xpmiss,omitempty"`
	FGM40To49    float64 `json:"fgm_40_49,omitempty"`
	FGM50Plus    float64 `json:"fgm_50p,omitempty"`
	FGMYds       float64 `json:"fgm_yds,omitempty"`
	FGMiss40To49 float64 `json:"fgmiss_40_49,omitempty"`
	FGMiss50Plus float64 `json:"fgmiss_50p,omitempty"`

	// Team defense / special teams (DEF)
	DefSack        float64 `json:"sack,omitempty"`
	DefInt         float64 `json:"int,omitempty"`
	DefFumRec      float64 `json:"fum_rec,omitempty"`
	DefFumTD       float64 `json:"def_fum_td,omitempty"`
	DefSafety      float64 `json:"safe,omitempty"`
	DefBlockKick   float64 `json:"blk_kick,omitempty"`
	DefKRTD        float64 `json:"def_kr_td,omitempty"`
	DefPRTD        float64 `json:"pr_td,omitempty"`
	DefPtsAllow0   float64 `json:"pts_allow_0,omitempty"`
	DefYdsAllow100 float64 `json:"yds_allow_0_100,omitempty"`

	// Individual defensive players (IDP)
	IDPTackle       float64 `json:"idp_tkl,omitempty"`
	IDPTackleSolo   float64 `json:"idp_tkl_solo,omitempty"`
	IDPTackleAssist float64 `json:"idp_tkl_ast,omitempty"`
	IDPSack         float64 `json:"idp_sack,omitempty"`
	IDPInt          float64 `json:"idp_int,omitempty"`
	IDPForcedFum    float64 `json:"idp_ff,omitempty"`
	IDPFumRec       float64 `json:"idp_fum_rec,omitempty"`
	IDPSafety       float64 `json:"idp_safe,omitempty"`
	IDPBlockKick    float64 `json:"idp_blk_kick,omitempty"`
}

// ADP returns the player's average draft position for the given format
// (ADPUndrafted if not being drafted, 0 if the stat was absent).
func (s ProjectionStats) ADP(adpType ADPType) float64 {
	switch adpType {
	case ADPTypeStandard:
		return s.ADPStandard
	case ADPTypeHalfPPR:
		return s.ADPHalfPPR
	case ADPTypePPR:
		return s.ADPPPR
	case ADPType2QB:
		return s.ADP2QB
	case ADPTypeDynasty:
		return s.ADPDynasty
	case ADPTypeDynastyStandard:
		return s.ADPDynastyStandard
	case ADPTypeDynastyHalfPPR:
		return s.ADPDynastyHalfPPR
	case ADPTypeDynastyPPR:
		return s.ADPDynastyPPR
	case ADPTypeDynasty2QB:
		return s.ADPDynasty2QB
	case ADPTypeIDP:
		return s.ADPIDP
	case ADPTypeIDP1QB:
		return s.ADPIDP1QB
	default:
		return 0
	}
}

// PlayerProjection is one entry from the projections endpoint: a player's
// projected season stats plus their average draft position in every format.
type PlayerProjection struct {
	PlayerID     string          `json:"player_id,omitempty"`
	Sport        string          `json:"sport,omitempty"`
	Season       string          `json:"season,omitempty"`
	SeasonType   string          `json:"season_type,omitempty"`
	Category     string          `json:"category,omitempty"`
	Company      string          `json:"company,omitempty"`
	GameID       string          `json:"game_id,omitempty"`
	Team         string          `json:"team,omitempty"`
	Stats        ProjectionStats `json:"stats,omitempty"`
	Player       *Player         `json:"player,omitempty"`
	UpdatedAt    int64           `json:"updated_at,omitempty"`
	LastModified int64           `json:"last_modified,omitempty"`

	// Null for season-long projections; populated in per-game entries.
	Week     int    `json:"week,omitempty"`
	Opponent string `json:"opponent,omitempty"`
	Status   string `json:"status,omitempty"`
	Date     string `json:"date,omitempty"`
}

type ProjectionOptions struct {
	Positions  []string // Filters on fantasy_positions (case-sensitive): "QB", "RB", "WR", "TE", "K", "DEF", or IDP "DL", "LB", "DB". Empty returns every position, including non-fantasy ones.
	SeasonType string   // "regular" (default), "pre", or "post"
	OrderBy    string   // any stat key, e.g. "adp_dynasty_2qb" or "pts_ppr"; controls server-side ordering only
}

// ListProjections retrieves player projections from a hidden Sleeper endpoint.
// This endpoint is not listed in the official Sleeper API documentation and may
// change or disappear without notice.
//
// Each entry's Stats holds projected season stats alongside the player's
// average draft position in every league format. To build a ranked ADP list,
// keep entries where Stats.ADP(adpType) < ADPUndrafted and sort ascending.
func (c *Client) ListProjections(ctx context.Context, sport Sport, season string, options ProjectionOptions) ([]PlayerProjection, error) {
	var errs []string
	if sport == "" {
		errs = append(errs, "sport is required")
	}
	if season == "" {
		errs = append(errs, "season is required")
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid request: %s", strings.Join(errs, "\n"))
	}

	if options.SeasonType == "" {
		options.SeasonType = "regular"
	}

	query := url.Values{}
	query.Set(queryParamSeasonType, options.SeasonType)
	if options.OrderBy != "" {
		query.Set(queryParamOrderBy, options.OrderBy)
	}
	for _, position := range options.Positions {
		query.Add(queryParamPosition, position)
	}
	endpoint := endpointHiddenBaseURL + fmt.Sprintf(endpointProjections, sport, season) + "?" + query.Encode()

	var projections []PlayerProjection
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting projections: %w", err)
	}

	if err := json.Unmarshal(by, &projections); err != nil {
		return nil, fmt.Errorf("unmarshaling projections: %w", err)
	}

	// Unknown sports return an empty array. Unknown seasons still return
	// player entries, just with empty Stats maps.
	if len(projections) == 0 {
		return nil, errors.New("projections not found: check that the sport is valid")
	}

	return projections, nil
}
