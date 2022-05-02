package models

import (
	"database/sql"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func Connect(user, password, name, host string) error {
	var db_url = "postgres://" + user + ":" + password + "@" + host + ":5432/" + name
	zero.Log.Debug().Msgf("connecting db to %s", db_url)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(db_url), pgdriver.WithInsecure(true)))
	db = bun.NewDB(sqldb, pgdialect.New())
	// defer db.Close()

	zero.Log.Info().Msgf("Connected to database: %v", db_url)

	//-- move the create table statements to another function.
	err := CreateSyllabusTable()
	return err
}
