package models

import (
	"context"
	"database/sql"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func InitDB(user, password, name, host string) (*bun.DB, error) {
	var url = "postgres://" + user + ":" + password + "@" + host + ":5432/" + name
	zero.Debugf("Connecting: %s", url)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url), pgdriver.WithInsecure(true)))
	db = bun.NewDB(sqldb, pgdialect.New())

	err := db.Ping()
	if err != nil {
		return db, err
	}

	zero.Infof("Connected: %v", url)
	err = SetupTables(false)
	return db, err
}

func SetupTables(reset bool) error {
	ctx := context.Background()
	if reset {
		db.NewDropTable().Model(&Syllabus{}).IfExists().Exec(ctx)
		db.NewDropTable().Model(&Contributor{}).IfExists().Exec(ctx)
		db.NewDropTable().Model(&Attachment{}).IfExists().Exec(ctx)
		db.NewDropTable().Model(&MagicToken{}).IfExists().Exec(ctx)
	}

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
