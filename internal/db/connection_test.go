package db

import (
	"context"
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
)

func TestConnect(t *testing.T) {

	time.Sleep(1 * time.Second)
	err := Connect("cosyl", "cosyl", "cosyl", "localhost")

	if err != nil {
		t.Error(err)
	}
}

func TestGetAllSyllabi(t *testing.T) {
	err := seedDB(t)
	if err != nil {
		t.Error(err)
	}

	syll, err := GetAllSyllabi()
	if err != nil {
		t.Error(err)
	}

	if len(syll) != 1 {
		t.Errorf("expected to have non-0 count of syllabi")
	}
}

func TestAddNewSyllabus(t *testing.T) {
	syll := models.Syllabus{Title: "Test Title 2", Description: "Test Description 2"}
	_, err := AddNewSyllabus(&syll)

	if err != nil {
		t.Error(err)
	}
}

func TestGetSyllabus(t *testing.T) {
	syll, err := GetSyllabus(1)
	if err != nil {
		t.Error(err)
	}

	if syll.Id != 1 {
		t.Errorf("Expecting ID to be 1, got %d", syll.Id)
	}
}

func TestUpdateSyllabus(t *testing.T) {
	syll := models.Syllabus{Title: "Test Title 1 (updated)", Description: "Test Description 1 (updated"}
	updated, err := UpdateSyllabus(1, &syll)

	if err != nil {
		t.Error(err)
	}

	if updated.Description != syll.Description {
		t.Error("Mismatch of updated syllabus")
	}
}

func TestDeleteSyllabus(t *testing.T) {
	err := DeleteSyllabus(1)
	if err != nil {
		t.Error(err)
	}
}

func seedDB(t *testing.T) error {
	Connect("test", "test", "test", "localhost")
	ctx := context.Background()

	db.NewDropTable().Model(&models.Syllabus{}).IfExists().Exec(ctx)
	_, err := db.NewCreateTable().Model((*models.Syllabus)(nil)).IfNotExists().Exec(ctx)

	if err != nil {
		return err
	}

	syll := models.Syllabus{Title: "Test Title 1", Description: "Test Description 1"}
	_, err = AddNewSyllabus(&syll)
	return err
}
