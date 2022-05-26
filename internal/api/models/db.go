package models

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *bun.DB

func InitDB(url string, fixturesDir string) (*bun.DB, error) {
	zero.Infof("connecting: %s", url) //-- todo this should not be logged
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

	err = runMigrations(url, sslMode)
	if err != nil {
		zero.Errorf("error running migrations: %v", err)
		log.Fatal(err)
	}

	if os.Getenv("TEST") == "true" {
		err = runFixtures(db, fixturesDir)
		if err != nil {
			zero.Errorf("error running fixtures: %v", err)
		}
	}
	return db, err
}

func runMigrations(url string, sslMode bool) error {
	if !sslMode {
		url = url + "?sslmode=disable"
	}

	migrationsDir := "file://internal/migrations"
	if os.Getenv("TEST") == "true" {
		migrationsDir = "file:///app/internal/migrations"
	}

	m, err := migrate.New(
		migrationsDir,
		url)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func runFixtures(db *bun.DB, dir string) error {

	db.RegisterModel(
		(*Syllabus)(nil),
		(*Attachment)(nil),
		(*Contributor)(nil),
		(*MagicToken)(nil),
	)

	fixture := dbfixture.New(db, dbfixture.WithTruncateTables())

	ctx := context.Background()
	err := fixture.Load(ctx, os.DirFS(dir), "syllabus.yml", "attachment.yml", "contributor.yml", "magic_token.yml")

	return err
}

func Shutdown() error {
	err := db.Close()
	return err
}
