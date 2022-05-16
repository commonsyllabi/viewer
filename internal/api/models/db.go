package models

import (
	"context"
	"database/sql"
	"strings"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func InitDB(url string) (*bun.DB, error) {
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

	zero.Infof("connected: %v", url) //should not be logged
	err = SetupTables(true)
	return db, err
}

func SetupTables(reset bool) error {
	ctx := context.Background()
	if reset {
		db.NewTruncateTable().Model(&Syllabus{}).Exec(ctx)
		db.NewTruncateTable().Model(&Contributor{}).Exec(ctx)
		db.NewTruncateTable().Model(&Attachment{}).Exec(ctx)
		db.NewTruncateTable().Model(&MagicToken{}).Exec(ctx)
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
