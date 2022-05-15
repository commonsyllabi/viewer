package models

import (
	"context"
	"io/ioutil"
	"testing"
	"time"
)

func TestInitDB(t *testing.T) {

	time.Sleep(1 * time.Second)
	// todo: pass this as an env variable for tests, connecting to the same docker-compose db hosts, but different postgres databases
	_, err := InitDB("postgres://cosyl:cosyl@localhost:5432/cosyl")

	if err != nil {
		t.Error(err)
	}
}

func mustSeedDB(t *testing.T) {
	InitDB("postgres://cosyl:cosyl@localhost:5432/cosyl")
	ctx := context.Background()

	//-- truncate table deletes all rows in a table
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
