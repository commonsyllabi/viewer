// Models keep track of data processing, such as returning slices for complex queries, or sanitizing input data
package models

import (
	"context"
)

type Syllabus struct {
	Id          int64 `bun:",pk,autoincrement"`
	Title       string
	Description string
	Contributor Contributor
}

func CreateSyllabusTable() error {
	ctx := context.Background()
	db.NewCreateTable().Model((*Syllabus)(nil)).IfNotExists().Exec(ctx)

	return nil
}

func GetAllSyllabi() ([]Syllabus, error) {
	ctx := context.Background()
	syllabi := make([]Syllabus, 0)

	err := db.NewSelect().Model(&syllabi).Scan(ctx, &syllabi)
	return syllabi, err
}

func AddNewSyllabus(syll *Syllabus) (Syllabus, error) {
	ctx := context.Background()

	_, err := db.NewInsert().Model(syll).Exec(ctx)
	return *syll, err
}

func UpdateSyllabus(id int, syll *Syllabus) (Syllabus, error) {
	ctx := context.Background()
	_, err := db.NewUpdate().Model(syll).WherePK().Exec(ctx)
	return *syll, err
}

func GetSyllabus(id int) (Syllabus, error) {
	ctx := context.Background()
	table := new(Syllabus)
	var syll Syllabus
	err := db.NewSelect().Model(table).Where("id = ?", id).Scan(ctx, &syll)
	return syll, err
}

func DeleteSyllabus(id int) error {
	ctx := context.Background()
	table := new(Syllabus)
	_, err := db.NewDelete().Model(table).Where("id = ?", id).Exec(ctx)

	return err
}
