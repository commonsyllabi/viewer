package models

import (
	"context"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {

	time.Sleep(1 * time.Second)
	err := Connect("cosyl", "cosyl", "cosyl", "localhost")

	if err != nil {
		t.Error(err)
	}
}

func seedDB(t *testing.T) error {
	Connect("test", "test", "test", "localhost")
	ctx := context.Background()

	db.NewDropTable().Model(&Syllabus{}).IfExists().Exec(ctx)
	_, err := db.NewCreateTable().Model((*Syllabus)(nil)).IfNotExists().Exec(ctx)

	if err != nil {
		return err
	}

	syll := Syllabus{Title: "Test Title 1", Description: "Test Description 1"}
	_, err = AddNewSyllabus(&syll)
	return err
}
