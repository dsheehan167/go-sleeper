# go-sleeper

A Go client library for the Sleeper Fantasy Football API.

## Features

- Full support for Sleeper's public API endpoints (leagues, users, drafts, players, transactions, rosters, avatars, and more)
- Rate limiting to avoid IP blocking (per Sleeper API guidelines)
- Structs and methods closely matching Sleeper's API documentation
- Simple, idiomatic Go interface
- Well-documented code and function comments

## Installation

```bash
go get github.com/dsheehan167/tacowire/go-sleeper
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"github.com/dsheehan167/tacowire/go-sleeper"
)

func main() {
	client, err := sleeper.NewClient(context.Background(), sleeper.Config{})
	if err != nil {
		panic(err)
	}

	// Example: Get a user by username or user_id
	user, err := client.GetUser(context.Background(), "username_or_userid")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}
```

## API Coverage

- Users: Get user, get user leagues
- Leagues: Get league, rosters, users, matchups, transactions, traded picks, brackets
- Drafts: Get draft, user drafts, league drafts, picks, traded picks
- Players: List all NFL players, trending players
- Rosters: Get roster details
- Transactions: Get transaction details
- Avatars: Download avatar images and thumbnails
- Sport State: Get current state for supported sports

## Rate Limiting

The client includes a rate limiter to help prevent exceeding Sleeper's API limits. Per Sleeper's documentation, you should stay under 1000 API calls per minute to avoid being IP-blocked.

## Trending Players Widget

To embed the official trending players list in your app or website, use the following HTML snippet provided by Sleeper:

```html
<iframe src="https://sleeper.app/embed/players/nfl/trending/add?lookback_hours=24&limit=25" width="350" height="500" allowtransparency="true" frameborder="0"></iframe>
```

## Documentation

- [Sleeper API Docs](https://docs.sleeper.com/)
- See GoDoc comments in the code for details on each method and struct.

## License

MIT
