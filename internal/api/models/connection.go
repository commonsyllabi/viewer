package models

import (
	"database/sql"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func InitDB(user, password, name, host string) error {
	var url = "postgres://" + user + ":" + password + "@" + host + ":5432/" + name
	zero.Log.Debug().Msgf("Connecting: %s", url)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url), pgdriver.WithInsecure(true)))
	db = bun.NewDB(sqldb, pgdialect.New())

	zero.Log.Info().Msgf("Connected: %v", url)
	err := CreateSyllabusTable()
	return err
}

func Shutdown() error {
	err := db.Close()
	return err
}
