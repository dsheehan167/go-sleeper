# go-sleeper

A Go client library for the [Sleeper](https://sleeper.com) fantasy sports API.

## Features

- Full coverage of Sleeper's public API (leagues, users, drafts, players, transactions, rosters, avatars, sport state)
- Built-in rate limiter to stay within Sleeper's API guidelines

## Installation

```bash
go get github.com/dsheehan167/go-sleeper
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	sleeper "github.com/dsheehan167/go-sleeper"
)

func main() {
	ctx := context.Background()

	client, err := sleeper.NewClient(ctx, sleeper.Config{})
	if err != nil {
		panic(err)
	}

	user, err := client.GetUser(ctx, "my_username")
	if err != nil {
		panic(err)
	}
	fmt.Println(user.DisplayName)

	leagues, err := client.GetUserLeagues(ctx, user.UserID, sleeper.SportNFL, "2024")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(leagues), "leagues")
}
```

## Configuration

All fields are optional — zero values use the defaults shown below.

```go
client, err := sleeper.NewClient(ctx, sleeper.Config{
	APIVersion:   sleeper.APIVersion1, // default: v1
	Timeout:      10 * time.Second,    // default: 30s
	RateLimitRPS: 10,                  // default: 15 requests/sec
	RateLimitBurst: 20,                // default: burst of 30
})
```

Per Sleeper's documentation, staying under 1000 requests per minute avoids IP blocks. The default rate limiter (15 RPS, burst 30) is well within that limit.

## API Coverage

| Area | Methods |
|------|---------|
| Users | `GetUser`, `GetUserLeagues` |
| Leagues | `GetLeague`, `GetLeagueRosters`, `GetLeagueUsers`, `GetLeagueMatchups`, `GetTransactions`, `GetLeagueTradedPicks`, `GetLeagueWinnersBracket`, `GetLeagueLosersBracket` |
| Drafts | `GetDraft`, `GetUserDrafts`, `GetLeagueDrafts`, `GetDraftPicks`, `GetDraftTradedPicks` |
| Players | `ListNFLPlayers`, `ListTrendingPlayers` |
| Avatars | `GetAvatarImage`, `GetAvatarThumbnail` |
| Sport State | `GetSportState` |

Supported sports: `SportNFL`, `SportNBA`, `SportMLB`, `SportNHL`.

## NFL Players

`ListNFLPlayers` returns the full player map (~5 MB). Sleeper recommends calling this **at most once per day** and caching the result on your own server rather than fetching it per-request.

## Trending Players Widget

To embed Sleeper's official trending list in a web page:

```html
<iframe src="https://sleeper.app/embed/players/nfl/trending/add?lookback_hours=24&limit=25" width="350" height="500" allowtransparency="true" frameborder="0"></iframe>
```

Please give attribution to Sleeper if you use their trending data.

## Documentation

- [Sleeper API docs](https://docs.sleeper.com/)
- GoDoc comments on every exported type and method

## License

MIT
