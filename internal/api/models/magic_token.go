package models

import "context"

type MagicToken struct {
	ID         int64  `bun:"id,pk,autoincrement"`
	Token      []byte `bun:"token"`
	SyllabusID int64
}

func CreateMagicTokenTable() error {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*MagicToken)(nil)).IfNotExists().Exec(ctx)

	return err
}

func GetTokenSyllabus(id int) (MagicToken, error) {
	ctx := context.Background()
	var token MagicToken
	err := db.NewSelect().Model(&token).Where("syllabus_id = ?", id).Scan(ctx)
	return token, err
}

func AddNewToken(token *MagicToken) (MagicToken, error) {
	ctx := context.Background()

	_, err := db.NewInsert().Model(token).Exec(ctx)
	return *token, err
}

func DeleteToken(id int) error {
	ctx := context.Background()
	table := new(MagicToken)
	_, err := db.NewDelete().Model(table).Where("id = ?", id).Exec(ctx)

	return err
}
