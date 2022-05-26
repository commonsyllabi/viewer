package models

import (
	"time"
)

type Contributor struct {
	CreatedAt time.Time   `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time   `bun:",nullzero,notnull,default:current_timestamp"`
	ID        int64       `bun:"id,pk,autoincrement"`
	Name      string      `bun:"name,notnull"`
	Email     string      `bun:"email,notnull,unique"`
	Syllabi   []*Syllabus `bun:"rel:has-many" form:"syllabi" json:"syllabi"`
}
