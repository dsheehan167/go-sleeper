package sleeper

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGetAvatarImage(t *testing.T) {
	// Load expected image data
	expectedData, err := os.ReadFile(filepath.Join("test", "avatar.png"))
	if err != nil {
		t.Fatalf("Failed to load test avatar image: %v", err)
	}

	tt := []struct {
		testcase     string
		avatarID     string
		expectedData []byte
		shouldPass   bool
	}{
		{
			"valid avatar ID",
			"cc12ec49965eb7856f84d71cf85306af",
			expectedData,
			true,
		},
		{
			"missing avatar ID",
			"",
			nil,
			false,
		},
	}

	ctx := context.Background()

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			data, err := testClient.GetAvatarImage(ctx, tc.avatarID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got response of length: %d", len(data))
				return
			}

			if !bytes.Equal(data, tc.expectedData) {
				t.Errorf("expected avatar data does not match retrieved data")
				return
			}

			t.Logf("retrieved avatar image of length: %d", len(data))
		})
	}
}
func TestGetAvatarThumbnail(t *testing.T) {
	// Load expected thumbnail data
	expectedData, err := os.ReadFile(filepath.Join("test", "avatarThumbnail.png"))
	if err != nil {
		t.Fatalf("Failed to load test avatar thumbnail: %v", err)
	}

	tt := []struct {
		testcase     string
		avatarID     string
		expectedData []byte
		shouldPass   bool
	}{
		{
			"valid avatar ID",
			"cc12ec49965eb7856f84d71cf85306af",
			expectedData,
			true,
		},
		{
			"missing avatar ID",
			"",
			nil,
			false,
		},
	}

	ctx := context.Background()

	for _, tc := range tt {
		t.Run(tc.testcase, func(t *testing.T) {
			data, err := testClient.GetAvatarThumbnail(ctx, tc.avatarID)
			if err != nil {
				if tc.shouldPass {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Logf("expected error: %v", err)
				return
			}

			if !tc.shouldPass {
				t.Errorf("expected failure but got response of length: %d", len(data))
				return
			}

			if !bytes.Equal(data, tc.expectedData) {
				t.Errorf("expected avatar thumbnail data does not match retrieved data")
				return
			}

			t.Logf("retrieved avatar thumbnail of length: %d", len(data))
		})
	}
}
