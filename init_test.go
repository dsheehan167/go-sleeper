package sleeper

import (
	"context"
	"log"
	"testing"
)

var testClient *Client

func TestMain(m *testing.M) {

	err := initTestClient()
	if err != nil {
		log.Fatalf("failed to initialize test client: %v", err)
	}

	m.Run()
}

func initTestClient() error {
	var err error
	testClient, err = NewClient(context.Background(), Config{})
	if err != nil {
		return err
	}

	return nil
}
