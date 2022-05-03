package models

import (
	"testing"
)

func TestGetAllSyllabi(t *testing.T) {
	mustSeedDB(t)

	syll, err := GetAllSyllabi()
	if err != nil {
		t.Error(err)
	}

	if len(syll) != 1 {
		t.Errorf("expected to have non-0 count of syllabi")
	}
}

// -- todo: handle when the response is empty
func TestAddNewSyllabus(t *testing.T) {
	syll := Syllabus{Title: "Test Title 2", Description: "Test Description 2"}
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

	if syll.ID != 1 {
		t.Errorf("Expecting ID to be 1, got %d", syll.ID)
	}
}

func TestUpdateSyllabus(t *testing.T) {
	mustSeedDB(t)

	syll := Syllabus{Title: "Test Title 1 (updated)", Description: "Test Description 1 (updated"}
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
