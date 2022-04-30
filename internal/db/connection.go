package db

import (
	"context"
	"database/sql"

	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func Connect(user, password, name, host string) error {
	var db_url = "postgres://" + user + ":" + password + "@" + host + ":5432/" + name
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(db_url), pgdriver.WithInsecure(true)))
	db = bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()
	ctx := context.Background()

	zero.Log.Info().Msgf("Connected to database: %v", db_url)

	_, err := db.NewCreateTable().Model((*models.Syllabus)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func AddNewSyllabus(syll models.Syllabus) (sql.Result, error) {
	ctx := context.Background()
	result, err := db.NewInsert().Model(syll).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	if err != nil {
		return result, err
	}
	return result, nil
}
