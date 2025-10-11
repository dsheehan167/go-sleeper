package sleeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// User represents a Sleeper user/league member.
type User struct {
	Avatar      string `json:"avatar,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	IsBot       bool   `json:"is_bot"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
}

// GetUser retrieves a User by identity (username or user_id of the user).
func (c *Client) GetUser(ctx context.Context, identity string) (*User, error) {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return nil, errors.New("userID is required")
	}

	endpoint := c.buildEndpoint(endpointUser, identity)

	var user *User
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	if err := json.Unmarshal(by, &user); err != nil {
		return nil, fmt.Errorf("unmarshaling user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserLeagues retrieves all leagues for a user is a member of for a given user id (must be the user_id, cannot be the username), sport, and season (season should be the year, e.g. 2025).
func (c *Client) GetUserLeagues(ctx context.Context, userID string, sport sport, season string) ([]*League, error) {

	var errs []string
	userID = strings.TrimSpace(userID)
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

	endpoint := c.buildEndpoint(endpointUserLeagues, userID, sport, season)

	var leagues []*League
	by, err := c.getRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("getting leagues for user: %w", err)
	}

	if err := json.Unmarshal(by, &leagues); err != nil {
		return nil, fmt.Errorf("unmarshaling leagues: %w", err)
	}

	return leagues, nil
}
