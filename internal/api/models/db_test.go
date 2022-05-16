package models

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var databaseTestURL string

func TestInitDB(t *testing.T) {

	time.Sleep(1 * time.Second)
	databaseTestURL = os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	_, err := InitDB(databaseTestURL)
	if err != nil {
		t.Error(err)
	}
}

func mustSeedDB(t *testing.T) {
	databaseTestURL = os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	InitDB(databaseTestURL)

	syll := Syllabus{Title: "Test Title 1", Description: "Test Description 1"}
	_, err := AddNewSyllabus(&syll)
	if err != nil {
		panic(err)
	}

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
		panic(err)
	}
}
