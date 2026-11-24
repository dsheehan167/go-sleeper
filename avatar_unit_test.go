package sleeper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

const unitAvatarBytes = "\x89PNG\r\n\x1a\nfake-image-data"

func TestUnitAvatar_GetImage(t *testing.T) {
	tt := []struct {
		testcase   string
		avatarID   string
		wantPath   string
		status     int
		body       string
		wantErr    string
		shouldPass bool
	}{
		{
			"valid avatar",
			"cc12ec49965eb7856f84d71cf85306af",
			"/avatars/cc12ec49965eb7856f84d71cf85306af",
			http.StatusOK,
			unitAvatarBytes,
			"",
			true,
		},
		{
			"avatar id is trimmed",
			" abc ",
			"/avatars/abc",
			http.StatusOK,
			unitAvatarBytes,
			"",
			true,
		},
		{
			"empty avatar id",
			"",
			"",
			0,
			"",
			"avatarID is required",
			false,
		},
		{
			"whitespace avatar id",
			"  ",
			"",
			0,
			"",
			"avatarID is required",
			false,
		},
		{
			"not found",
			"missing",
			"/avatars/missing",
			http.StatusNotFound,
			``,
			"API error: 404",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respondExact(t, "sleepercdn.com", tc.wantPath, "", tc.status, tc.body))

			image, err := client.GetAvatarImage(context.Background(), tc.avatarID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("expected error containing %q, got %q", tc.wantErr, err)
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d bytes", len(image))
				return
			}

			if string(image) != tc.body {
				t.Errorf("expected image bytes %q, got %q", tc.body, image)
				return
			}
		})
	}
}

func TestUnitAvatar_GetThumbnail(t *testing.T) {
	tt := []struct {
		testcase   string
		avatarID   string
		wantPath   string
		status     int
		body       string
		wantErr    string
		shouldPass bool
	}{
		{
			"valid thumbnail",
			"cc12ec49965eb7856f84d71cf85306af",
			"/avatars/thumbs/cc12ec49965eb7856f84d71cf85306af",
			http.StatusOK,
			unitAvatarBytes,
			"",
			true,
		},
		{
			"avatar id is trimmed",
			"\tabc\n",
			"/avatars/thumbs/abc",
			http.StatusOK,
			unitAvatarBytes,
			"",
			true,
		},
		{
			"empty avatar id",
			"",
			"",
			0,
			"",
			"avatarID is required",
			false,
		},
		{
			"server error",
			"abc",
			"/avatars/thumbs/abc",
			http.StatusInternalServerError,
			``,
			"API error: 500",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			client := newUnitTestClient(t, Config{}, respondExact(t, "sleepercdn.com", tc.wantPath, "", tc.status, tc.body))

			thumbnail, err := client.GetAvatarThumbnail(context.Background(), tc.avatarID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("expected error containing %q, got %q", tc.wantErr, err)
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got %d bytes", len(thumbnail))
				return
			}

			if string(thumbnail) != tc.body {
				t.Errorf("expected thumbnail bytes %q, got %q", tc.body, thumbnail)
				return
			}
		})
	}
}
