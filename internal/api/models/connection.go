package models

import (
	"context"
	"database/sql"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func InitDB(user, password, name, host string) error {
	var url = "postgres://" + user + ":" + password + "@" + host + ":5432/" + name
	zero.Log.Debug().Msgf("Connecting: %s", url)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url), pgdriver.WithInsecure(true)))
	DB = bun.NewDB(sqldb, pgdialect.New())

	err := DB.Ping()
	if err != nil {
		return err
	}

	zero.Log.Info().Msgf("Connected: %v", url)
	err = setupTables(true)
	return err
}

func setupTables(reset bool) error {
	ctx := context.Background()
	if reset {
		DB.NewDropTable().Model(&Syllabus{}).IfExists().Exec(ctx)
		DB.NewDropTable().Model(&Contributor{}).IfExists().Exec(ctx)
		DB.NewDropTable().Model(&Attachment{}).IfExists().Exec(ctx)
	}

	if err := CreateSyllabiTable(); err != nil {
		return err
	}
	if err := CreateContributorsTable(); err != nil {
		return err
	}
	err := CreateAttachmentsTable()

	return err
}

func Shutdown() error {
	err := DB.Close()
	return err
}
