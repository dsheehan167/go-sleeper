package sleeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type sport string

const (
	SportNFL sport = "nfl"
	SportMLB sport = "mlb"
	SportNBA sport = "nba"
	SportNHL sport = "nhl"
)

type SportState struct {
	Week               int            `json:"week,omitempty"`              // week
	SeasonType         string         `json:"season_type,omitempty"`       // pre, post, regular
	SeasonStartDate    string         `json:"season_start_date,omitempty"` // regular season start
	Season             string         `json:"season,omitempty"`            // current season
	PreviousSeason     FlexibleString `json:"previous_season,omitempty"`
	Leg                int            `json:"leg,omitempty"`                  // week of regular season
	LeagueSeason       string         `json:"league_season,omitempty"`        // active season for leagues
	LeagueCreateSeason string         `json:"league_create_season,omitempty"` // flips in December
	DisplayWeek        int            `json:"display_week,omitempty"`         // Which week to display in UI, can be different than week
}

func (c *Client) GetSportState(ctx context.Context, sport sport) (*SportState, error) {
	if sport == "" {
		return nil, errors.New("sport is required")
	}

	endpoint := c.buildEndpoint(endpointSportState, string(sport))

	var state *SportState
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting sport state: %w", err)
	}

	if err := json.Unmarshal(by, &state); err != nil {
		return nil, fmt.Errorf("unmarshaling sport state: %w", err)
	}

	if state == nil {
		return nil, errors.New("sport state not found")
	}

	return state, nil
}
