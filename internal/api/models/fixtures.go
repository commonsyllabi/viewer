package models

import (
	"context"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

func RunFixtures(db *bun.DB, dir string) error {

	db.RegisterModel(
		(*Syllabus)(nil),
		(*Attachment)(nil),
		(*Contributor)(nil),
		(*MagicToken)(nil),
	)

	fixture := dbfixture.New(db, dbfixture.WithTruncateTables())

	ctx := context.Background()
	err := fixture.Load(ctx, os.DirFS(dir), "syllabus.yml", "attachment.yml", "contributor.yml", "magic_token.yml")

	return err
}
