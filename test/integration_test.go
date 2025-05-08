package main

import (
	"os"
	"testing"

	"github.com/Shriharsh07/chaintrack/config"
)

func TestDBConnection(t *testing.T) {
	if os.Getenv("DB_HOST") == "" {
		t.Skip("Skipping integration test")
	}

	// Example test
	err := config.ConnectDB()
	if err != nil {
		t.Fatalf("DB connection failed: %v", err)
	}
}
