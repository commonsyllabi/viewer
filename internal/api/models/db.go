package models

import (
	"database/sql"
	"strings"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func InitDB(url string, fixturesDir string) (*bun.DB, error) {
	zero.Debugf("connecting: %s", url) //-- todo this should not be logged
	sslMode := false
	if strings.HasSuffix(url, "sslmode=require") {
		sslMode = true
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url), pgdriver.WithInsecure(!sslMode)))

	db = bun.NewDB(sqldb, pgdialect.New())

	err := db.Ping()
	if err != nil {
		return db, err
	}

	err = SetupTables()
	if err != nil {
		zero.Errorf("error setting up tables: %v", err)
	}
	err = RunFixtures(db, fixturesDir) //-- truncates tables
	if err != nil {
		zero.Errorf("error running fixtures: %v", err)
	}
	return db, err
}

// SetupTable creates all tables in the database. This should actually be RunFixtures
func SetupTables() error {
	if err := CreateSyllabiTable(); err != nil {
		return err
	}
	if err := CreateContributorsTable(); err != nil {
		return err
	}
	if err := CreateMagicTokenTable(); err != nil {
		return err
	}
	err := CreateAttachmentsTable()

	return err
}

func Shutdown() error {
	err := db.Close()
	return err
}
