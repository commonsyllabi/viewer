package models

import (
	"context"
	"io/ioutil"
	"testing"
	"time"
)

func TestInitDB(t *testing.T) {

	time.Sleep(1 * time.Second)
	err := InitDB("cosyl", "cosyl", "cosyl", "localhost")

	if err != nil {
		t.Error(err)
	}
}

func mustSeedDB(t *testing.T) {
	InitDB("test", "test", "test", "localhost")
	ctx := context.Background()

	db.NewDropTable().Model(&Syllabus{}).IfExists().Exec(ctx)
	_, err := db.NewCreateTable().Model((*Syllabus)(nil)).IfNotExists().Exec(ctx)

	if err != nil {
		panic(err)
	}

	syll := Syllabus{Title: "Test Title 1", Description: "Test Description 1"}
	_, err = AddNewSyllabus(&syll)
	if err != nil {
		panic(err)
	}

	db.NewDropTable().Model(&Attachment{}).IfExists().Exec(ctx)
	_, err = db.NewCreateTable().Model((*Attachment)(nil)).IfNotExists().Exec(ctx)

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
