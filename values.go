package sleeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ScoringFormat selects the league format used to compute player values.
type ScoringFormat string

const (
	ScoringFormatStandard        ScoringFormat = "std"
	ScoringFormatHalfPPR         ScoringFormat = "half_ppr"
	ScoringFormatPPR             ScoringFormat = "ppr"
	ScoringFormat2QB             ScoringFormat = "2qb"
	ScoringFormatDynastyStandard ScoringFormat = "dynasty_std"
	ScoringFormatDynastyHalfPPR  ScoringFormat = "dynasty_half_ppr"
	ScoringFormatDynastyPPR      ScoringFormat = "dynasty_ppr"
	ScoringFormatDynasty2QB      ScoringFormat = "dynasty_2qb"
	ScoringFormatIDP             ScoringFormat = "idp"
)

type PlayerValuesOptions struct {
	IncludeIDP bool // Include defensive (IDP) players in the response
	IsDynasty  bool // Sent as is_dynasty to match the Sleeper app, but the server ignores it; dynasty values are selected via the dynasty_* formats
}

// GetPlayerValues retrieves player values from a hidden Sleeper endpoint. This
// endpoint is not listed in the official Sleeper API documentation and may
// change or disappear without notice.
//
// The response maps player IDs to a value score where a higher value indicates
// a more valuable player. This is a value model, not draft position: values
// correlate with ADP but do not match the ADP shown in the Sleeper draft lobby.
// For actual average draft positions, use the adp_* stats from ListProjections.
//
// The format is the only scoring dimension the endpoint supports; there are no
// finer-grained variants (TEP, superflex, best ball) — unknown formats and query
// parameters are silently ignored and return an empty result. NFL supports all
// scoring formats; other sports may only return data for ScoringFormatStandard
// and the current season.
func (c *Client) GetPlayerValues(ctx context.Context, sport Sport, season string, format ScoringFormat, options PlayerValuesOptions) (map[string]float64, error) {
	var errs []string
	if sport == "" {
		errs = append(errs, "sport is required")
	}
	if season == "" {
		errs = append(errs, "season is required")
	}
	if format == "" {
		errs = append(errs, "format is required")
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid request: %s", strings.Join(errs, "\n"))
	}

	endpoint := endpointHiddenBaseURL + fmt.Sprintf(endpointPlayerValues, sport, season, format)
	endpoint = addQueryParams(endpoint, map[string]string{
		queryParamIDP:       strconv.FormatBool(options.IncludeIDP),
		queryParamIsDynasty: strconv.FormatBool(options.IsDynasty),
	})

	var values map[string]float64
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting player values: %w", err)
	}

	if err := json.Unmarshal(by, &values); err != nil {
		return nil, fmt.Errorf("unmarshaling player values: %w", err)
	}

	// The endpoint returns an empty object (not an error) for unknown
	// seasons or formats.
	if len(values) == 0 {
		return nil, errors.New("player values not found: check that the season and format are valid")
	}

	return values, nil
}
