package sleeper

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	avatarBaseURL      = "https://sleepercdn.com/avatars"
	avatarThumbnailURL = "https://sleepercdn.com/avatars/thumbs"
)

// GetAvatarImage downloads the full-size avatar image and returns the raw bytes
func (c *Client) GetAvatarImage(ctx context.Context, avatarID string) ([]byte, error) {
	avatarID = strings.TrimSpace(avatarID)
	if avatarID == "" {
		return nil, errors.New("avatarID is required")
	}

	endpoint := fmt.Sprintf("%s/%s", avatarBaseURL, avatarID)
	return c.getRequest(ctx, endpoint)
}

// GetAvatarThumbnail downloads the thumbnail avatar image and returns the raw bytes
func (c *Client) GetAvatarThumbnail(ctx context.Context, avatarID string) ([]byte, error) {
	avatarID = strings.TrimSpace(avatarID)
	if avatarID == "" {
		return nil, errors.New("avatarID is required")
	}

	endpoint := fmt.Sprintf("%s/%s", avatarThumbnailURL, avatarID)
	return c.getRequest(ctx, endpoint)
}
