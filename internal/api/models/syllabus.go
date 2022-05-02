// Models keep track of data processing, such as returning slices for complex queries, or sanitizing input data
package models

type Syllabus struct {
	Id          int64 `bun:",pk,autoincrement"`
	Title       string
	Description string
	Contributor Contributor
}
