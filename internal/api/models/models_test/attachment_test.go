package models_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
)

const singleTestFile = "../../../../pkg/commoncartridge/test_files/test_01.imscc"

func TestGetAllAttachments(t *testing.T) {
	mustSeedDB(t)

	att, err := models.GetAllAttachments()
	if err != nil {
		t.Error(err)
	}

	if len(att) != 1 {
		t.Errorf("expected to have non-0 count of attachments")
	}
}

func TestAddNewAttachment(t *testing.T) {
	bytes, err := ioutil.ReadFile(singleTestFile)
	if err != nil {
		t.Error(err)
	}
	att := models.Attachment{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "test_01.imscc",
		File:      bytes,
		Type:      "zip",
	}

	_, err = models.AddNewAttachment(&att)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAttachment(t *testing.T) {
	syll, err := models.GetAttachment(1)
	if err != nil {
		t.Error(err)
	}

	if syll.ID != 1 {
		t.Errorf("Expecting ID to be 1, got %d", syll.ID)
	}
}

func TestUpdateAttachment(t *testing.T) {
	mustSeedDB(t)

	bytes, err := ioutil.ReadFile(singleTestFile)
	if err != nil {
		t.Error(err)
	}
	att := models.Attachment{
		UpdatedAt: time.Now(),
		Name:      "test_01.imscc (updated)",
		File:      bytes,
		Type:      "zip",
	}

	updated, err := models.UpdateAttachment(1, &att)

	if err != nil {
		t.Error(err)
	}

	if updated.Name != att.Name {
		t.Error("Mismatch of updated Attachment")
	}
}

func TestDeleteAttachment(t *testing.T) {
	err := models.DeleteAttachment(1)
	if err != nil {
		t.Error(err)
	}
}
