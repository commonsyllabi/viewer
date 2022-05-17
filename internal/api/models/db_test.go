package models

import (
	"os"
	"testing"
	"time"
)

var databaseTestURL string

func TestInitDB(t *testing.T) {

	time.Sleep(1 * time.Second)
	databaseTestURL = os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	_, err := InitDB(databaseTestURL, "../models/fixtures")
	if err != nil {
		t.Error(err)
	}
}

func mustSeedDB(t *testing.T) {
	databaseTestURL = os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	_, err := InitDB(databaseTestURL, "../models/fixtures")
	if err != nil {
		panic(err)
	}
}
