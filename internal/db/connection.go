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
	zero.Log.Debug().Msgf("connecting db to %s", db_url)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(db_url), pgdriver.WithInsecure(true)))
	db = bun.NewDB(sqldb, pgdialect.New())
	// defer db.Close()
	ctx := context.Background()

	zero.Log.Info().Msgf("Connected to database: %v", db_url)

	_, err := db.NewCreateTable().Model((*models.Syllabus)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetAllSyllabi() ([]models.Syllabus, error) {
	ctx := context.Background()
	syllabi := make([]models.Syllabus, 0)

	err := db.NewSelect().Model(&syllabi).Scan(ctx, &syllabi)
	return syllabi, err
}

func AddNewSyllabus(syll *models.Syllabus) (models.Syllabus, error) {
	ctx := context.Background()

	_, err := db.NewInsert().Model(syll).Exec(ctx)
	return *syll, err
}

func UpdateSyllabus(id int, syll *models.Syllabus) (models.Syllabus, error) {
	ctx := context.Background()
	_, err := db.NewUpdate().Model(syll).WherePK().Exec(ctx)
	if err != nil {
		return *syll, err
	}
	return *syll, nil
}

func GetSyllabus(id int) (models.Syllabus, error) {
	ctx := context.Background()
	table := new(models.Syllabus)
	var syll models.Syllabus
	err := db.NewSelect().Model(table).Where("id = ?", id).Scan(ctx, &syll)
	return syll, err
}

func DeleteSyllabus(id int) error {
	ctx := context.Background()
	table := new(models.Syllabus)
	_, err := db.NewDelete().Model(table).Where("id = ?", id).Exec(ctx)

	return err
}
