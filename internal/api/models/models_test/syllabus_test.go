package models_test

import (
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
)

//-- to not recreate the whole table everytime, one can also do all transactions, and then rollback rather than commit at the end of each test
func TestGetAllSyllabi(t *testing.T) {
	mustSeedDB(t)

	syll, err := models.GetAllSyllabi()
	if err != nil {
		t.Error(err)
	}

	if len(syll) == 0 {
		t.Errorf("expected to have non-0 count of syllabi, got %d", len(syll))
	}
}

// -- todo: handle when the response is empty
func TestAddNewSyllabus(t *testing.T) {
	syll := models.Syllabus{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       "Test Title 2",
		Description: "Test Description 2",
	}
	_, err := models.AddNewSyllabus(&syll)

	if err != nil {
		t.Error(err)
	}
}

func TestGetSyllabus(t *testing.T) {
	syll, err := models.GetSyllabus(1)
	if err != nil {
		t.Error(err)
	}

	if syll.ID != 1 {
		t.Errorf("Expecting ID to be 1, got %d", syll.ID)
	}
}

func TestUpdateSyllabus(t *testing.T) {
	mustSeedDB(t)

	syll := models.Syllabus{
		UpdatedAt:   time.Now(),
		Title:       "Test Title 1 (updated)",
		Description: "Test Description 1 (updated",
	}
	updated, err := models.UpdateSyllabus(1, &syll)

	if err != nil {
		t.Error(err)
	}

	if updated.Description != syll.Description {
		t.Error("Mismatch of updated syllabus")
	}
}

func TestDeleteSyllabus(t *testing.T) {
	err := models.DeleteSyllabus(1)
	if err != nil {
		t.Error(err)
	}
}
