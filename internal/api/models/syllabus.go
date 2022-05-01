package models

type Syllabus struct {
	Id          int64 `bun:",pk,autoincrement"`
	Title       string
	Description string
	Contributor Contributor
}
