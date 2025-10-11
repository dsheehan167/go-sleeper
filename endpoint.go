package sleeper

import (
	"fmt"
	"net/url"
)

const (
	endpointBaseURL = "https://api.sleeper.app"

	endpointSportState = "/state/%s"

	endpointUser        = "/user/%s"
	endpointUserLeagues = "/user/%s/leagues/%s/%s"

	endpointLeague               = "/league/%s"
	endpointLeagueRosters        = "/league/%s/rosters"
	endpointLeagueUsers          = "/league/%s/users"
	endpointLeagueMatchups       = "/league/%s/matchups/%d"
	endpointLeagueTransactions   = "/league/%s/transactions/%d"
	endpointLeagueTradedPicks    = "/league/%s/traded_picks"
	endpointLeagueWinnersBracket = "/league/%s/winners_bracket"
	endpointLeagueLosersBracket  = "/league/%s/losers_bracket"

	endpointNFLPlayers      = "/players/nfl"
	endpointTrendingPlayers = "/players/%s/trending/%s"

	endpointDraft            = "/draft/%s"
	endpointUserDrafts       = "/user/%s/drafts/%s/%s"
	endpointLeagueDrafts     = "/league/%s/drafts"
	endpointDraftTradedPicks = "/draft/%s/traded_picks"
	endpointDraftPicks       = "/draft/%s/picks"

	// Query parameters
	queryParamLookbackHours = "lookback_hours"
	queryParamLimit         = "limit"
)

func (c *Client) buildEndpoint(template string, args ...interface{}) string {
	return c.baseURL + fmt.Sprintf(template, args...)
}

// addQueryParams adds query parameters to an endpoint URL
func addQueryParams(endpoint string, params map[string]string) string {
	if len(params) == 0 {
		return endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return endpoint
	}

	q := u.Query()
	for key, value := range params {
		if value != "" {
			q.Add(key, value)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}
