package models

type Contributor struct {
	Id      int64 `bun:",pk,autoincrement"`
	name    string
	email   string
	syllabi []Syllabus
}
