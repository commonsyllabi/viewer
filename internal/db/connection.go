package db

import (
	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func Connect(user, password, name, host string) error {
	var db_url = "postgres://" + user + ":" + password + "@" + host + ":5432/" + name
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: name,
	})

	defer db.Close()

	zero.Log.Info().Msgf("Connected to database: %v", db_url)

	syll := models.Syllabus{Id: 0, Title: "Others", Description: "They will always be there."}
	err := db.Model(&syll).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func AddNewSyllabus(syll models.Syllabus) error {
	return nil
}
