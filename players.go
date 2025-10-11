package sleeper

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"encoding/json"
	"fmt"
)

type trendingType string

const (
	TrendingTypeAdd  trendingType = "add"
	TrendingTypeDrop trendingType = "drop"
)

// Player represents a player in the Sleeper system
type Player struct {
	PlayerID              string          `json:"player_id"`
	FirstName             string          `json:"first_name"`
	LastName              string          `json:"last_name"`
	FullName              string          `json:"full_name"`
	Status                string          `json:"status"`
	Sport                 string          `json:"sport"`
	Position              string          `json:"position"`
	Team                  *string         `json:"team"`
	TeamAbbr              *string         `json:"team_abbr"`
	Number                *int            `json:"number"`
	Age                   *int            `json:"age"`
	Height                string          `json:"height"`
	Weight                string          `json:"weight"`
	College               string          `json:"college"`
	HighSchool            *string         `json:"high_school"`
	BirthDate             *string         `json:"birth_date"`
	BirthCity             *string         `json:"birth_city"`
	BirthState            *string         `json:"birth_state"`
	BirthCountry          *string         `json:"birth_country"`
	YearsExp              int             `json:"years_exp"`
	Active                bool            `json:"active"`
	SearchRank            int             `json:"search_rank"`
	SearchFirstName       string          `json:"search_first_name"`
	SearchLastName        string          `json:"search_last_name"`
	SearchFullName        string          `json:"search_full_name"`
	FantasyPositions      []string        `json:"fantasy_positions"`
	DepthChartPosition    *string         `json:"depth_chart_position"`
	DepthChartOrder       *int            `json:"depth_chart_order"`
	InjuryStatus          *string         `json:"injury_status"`
	InjuryBodyPart        *string         `json:"injury_body_part"`
	InjuryNotes           *string         `json:"injury_notes"`
	InjuryStartDate       *string         `json:"injury_start_date"`
	PracticeParticipation *string         `json:"practice_participation"`
	PracticeDescription   *string         `json:"practice_description"`
	NewsUpdated           *int64          `json:"news_updated"`
	TeamChangedAt         *string         `json:"team_changed_at"`
	Hashtag               string          `json:"hashtag"`
	Metadata              *PlayerMetadata `json:"metadata"`

	// Third-party IDs
	ESPNID        *int    `json:"espn_id"`
	YahooID       *int    `json:"yahoo_id"`
	RotowireID    *int    `json:"rotowire_id"`
	RotoworldID   *int    `json:"rotoworld_id"`
	GSIID         *string `json:"gsis_id"`
	SportradarID  string  `json:"sportradar_id"`
	StatsID       *int    `json:"stats_id"`
	FantasyDataID *int    `json:"fantasy_data_id"`
	SwishID       *int    `json:"swish_id"`
	OptaID        *string `json:"opta_id"`
	PandascoreID  *string `json:"pandascore_id"`
	OddsjamID     *string `json:"oddsjam_id"`
	KalshiID      *string `json:"kalshi_id"`
}

type PlayerMetadata struct {
	ChannelID  string `json:"channel_id"`
	RookieYear string `json:"rookie_year"`
}

type TrendingPlayerOptions struct {
	LookbackHours int `json:"lookback_hours,omitempty"` // Number of hours to look back for trending players (default: 24)
	Limit         int `json:"limit,omitempty"`          // Maximum number of players to return (default: 10, max: 50)
}

// ListNFLPlayers retrieves the full map of NFL players from the Sleeper API.
// This endpoint returns a large payload (~5MB) and should only be called once per day
// to update your local cache of player IDs and metadata. The response maps player IDs
// (e.g., "1042", "2403", "CAR") to player information, which is necessary for resolving
// player references in rosters and draft picks. Do not call this endpoint on every lookup;
// instead, store the results on your own server and refresh them daily at most.
func (c *Client) ListNFLPlayers(ctx context.Context) (map[string]Player, error) {
	endpoint := c.buildEndpoint(endpointNFLPlayers)

	var players map[string]Player
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting players: %w", err)
	}

	if err := json.Unmarshal(by, &players); err != nil {
		return nil, fmt.Errorf("unmarshaling players: %w", err)
	}

	if players == nil {
		return nil, errors.New("players not found")
	}

	return players, nil
}

// ListTrendingPlayers retrieves a list of trending players based on adds or drops in the past 24 hours from the Sleeper API.
// Please give attribution to Sleeper if you use their trending data. This endpoint can be used to get trending players
// for a given sport and type ("add" or "drop") with optional lookback_hours and limit parameters.
//
// If you wish to embed the official trending list in your app or website, use the following HTML snippet provided by Sleeper:
//
// <iframe src="https://sleeper.app/embed/players/nfl/trending/add?lookback_hours=24&limit=25" width="350" height="500" allowtransparency="true" frameborder="0"></iframe>
func (c *Client) ListTrendingPlayers(ctx context.Context, sport sport, trendingType trendingType, options TrendingPlayerOptions) ([]*Player, error) {
	var errs []string
	if sport == "" {
		errs = append(errs, "sport is required")
	}
	if trendingType == "" {
		errs = append(errs, "trendingType is required")
	}
	if err := options.validate(); err != nil {
		errs = append(errs, fmt.Sprintf("invalid options: %v", err))
	}
	if len(errs) > 0 {
		return nil, errors.New("invalid request:\n" + fmt.Sprintf("%s", errs))
	}

	endpoint := c.buildEndpoint(endpointTrendingPlayers, sport, trendingType)
	queryParms := make(map[string]string)
	if options.LookbackHours != 0 {
		queryParms[queryParamLookbackHours] = strconv.Itoa(options.LookbackHours)
	}
	if options.Limit != 0 {
		queryParms[queryParamLimit] = strconv.Itoa(options.Limit)
	}
	if len(queryParms) > 0 {
		endpoint = addQueryParams(endpoint, queryParms)
	}

	var players []*Player
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting trending players: %w", err)
	}

	if err := json.Unmarshal(by, &players); err != nil {
		return nil, fmt.Errorf("unmarshaling trending players: %w", err)
	}

	if len(players) == 0 {
		return nil, errors.New("trending players not found")
	}

	return players, nil
}

func (o *TrendingPlayerOptions) validate() error {
	var errs []string
	if o.LookbackHours < 0 {
		errs = append(errs, "LookbackHours must be greater than or equal to zero")
	}
	if o.Limit < 0 {
		errs = append(errs, "Limit must be greater than or equal to zero")
	}
	if o.Limit > 50 {
		errs = append(errs, "Limit cannot exceed 50")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
