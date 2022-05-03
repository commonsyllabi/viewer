package models

import (
	"io/ioutil"
	"testing"
)

const singleTestFile = "../../../pkg/commoncartridge/test_files/test_01.imscc"

func TestGetAllAttachments(t *testing.T) {
	mustSeedDB(t)

	att, err := GetAllAttachments()
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
	att := Attachment{
		Name: "test_01.imscc",
		File: bytes,
		Type: "zip",
	}

	_, err = AddNewAttachment(&att)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAttachment(t *testing.T) {
	syll, err := GetAttachment(1)
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
	att := Attachment{
		Name: "test_01.imscc (updated)",
		File: bytes,
		Type: "zip",
	}

	updated, err := UpdateAttachment(1, &att)

	if err != nil {
		t.Error(err)
	}

	if updated.Name != att.Name {
		t.Error("Mismatch of updated Attachment")
	}
}

func TestDeleteAttachment(t *testing.T) {
	err := DeleteAttachment(1)
	if err != nil {
		t.Error(err)
	}
}
