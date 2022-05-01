package models

type Contributor struct {
	Id      int64 `bun:",pk,autoincrement"`
	Name    string
	Email   string
	Syllabi []Syllabus
}
