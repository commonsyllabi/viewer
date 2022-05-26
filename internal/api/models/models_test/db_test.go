package models_test

import (
	"os"
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
)

var databaseTestURL string

func TestInitDB(t *testing.T) {

	time.Sleep(1 * time.Second)
	databaseTestURL = os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	_, err := models.InitDB(databaseTestURL, "../../models/fixtures")
	if err != nil {
		t.Error(err)
	}
}

var isDBSeeded = false

//-- todo here we should check whether the db is already initialized or not
func mustSeedDB(t *testing.T) {
	if !isDBSeeded {
		databaseTestURL = os.Getenv("DATABASE_TEST_URL")
		if databaseTestURL == "" {
			databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
		}
		_, err := models.InitDB(databaseTestURL, "../../models/fixtures")
		if err != nil {
			panic(err)
		}

		isDBSeeded = true
	}
}
