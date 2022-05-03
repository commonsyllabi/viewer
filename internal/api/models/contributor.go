package models

import "context"

type Contributor struct {
	ID      int64 `bun:"id,pk,autoincrement"`
	Name    string
	Email   string
	Syllabi []*Syllabus `bun:"rel:has-many" form:"syllabi" json:"syllabi"`
}

func CreateContributorsTable() error {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*Contributor)(nil)).IfNotExists().Exec(ctx)

	return err
}
