// Models keep track of data processing, such as returning slices for complex queries, or sanitizing input data
package models

import (
	"context"
)

type Syllabus struct {
	ID          int64        `bun:"id,pk,autoincrement"`
	Title       string       `form:"title" json:"title"`
	Description string       `form:"description" json:"description"`
	Attachments []Attachment `bun:"rel:has-many"`
}

func CreateSyllabiTable() error {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*Syllabus)(nil)).IfNotExists().Exec(ctx)

	return err
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
	var syll Syllabus
	err := db.NewSelect().Model(&syll).Where("id = ?", id).Relation("Attachments").Scan(ctx)
	return syll, err
}

func DeleteSyllabus(id int) error {
	ctx := context.Background()
	table := new(Syllabus)
	_, err := db.NewDelete().Model(table).Where("id = ?", id).Exec(ctx) //-- what to do with dangling attachments?

	return err
}
