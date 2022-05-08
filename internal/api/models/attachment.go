package models

import "context"

type Attachment struct {
	ID         int64  `bun:"id,pk,autoincrement"`
	Name       string `form:"name"`
	File       []byte
	Type       string
	SyllabusID int64
	Syllabus   Syllabus `bun:"belongs-to,join:syllabus_id=id"`
	//cartridge  commoncartridge.Cartridge //-- have this directly? or leave it as file interface?
}

func CreateAttachmentsTable() error {
	ctx := context.Background()
	_, err := DB.NewCreateTable().Model((*Attachment)(nil)).IfNotExists().Exec(ctx)

	return err
}

func GetAllAttachments() ([]Attachment, error) {
	ctx := context.Background()
	att := make([]Attachment, 0)

	err := DB.NewSelect().Model(&att).Scan(ctx, &att)
	return att, err
}

func AddNewAttachment(att *Attachment) (Attachment, error) {
	ctx := context.Background()

	_, err := DB.NewInsert().Model(att).Exec(ctx)
	return *att, err
}

func UpdateAttachment(id int, att *Attachment) (Attachment, error) {
	ctx := context.Background()
	_, err := DB.NewUpdate().Model(att).WherePK().Exec(ctx)
	return *att, err
}

func GetAttachment(id int) (Attachment, error) {
	ctx := context.Background()
	table := new(Attachment)
	var att Attachment
	err := DB.NewSelect().Model(table).Where("id = ?", id).Scan(ctx, &att)
	return att, err
}

func DeleteAttachment(id int) error {
	ctx := context.Background()
	table := new(Attachment)
	_, err := DB.NewDelete().Model(table).Where("id = ?", id).Exec(ctx)

	return err
}
