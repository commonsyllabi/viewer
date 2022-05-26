package models

import (
	"context"
	"time"
)

type Attachment struct {
	CreatedAt          time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt          time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	ID                 int64     `bun:"id,pk,autoincrement"`
	Name               string    `bun:"name,notnull" form:"name"`
	File               []byte    `bun:"file,notnull"`
	Type               string    `bun:"type,notnull"`
	SyllabusAttachedID int64     `yaml:"syllabus_attached_id"`
	Syllabus           *Syllabus `bun:"rel:belongs-to,join:syllabus_attached_id=id"`
}

func CreateAttachmentsTable() error {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*Attachment)(nil)).IfNotExists().Exec(ctx)

	return err
}

func GetAllAttachments() ([]Attachment, error) {
	ctx := context.Background()
	att := make([]Attachment, 0)

	err := db.NewSelect().Model(&att).Scan(ctx, &att)
	return att, err
}

func AddNewAttachment(att *Attachment) (Attachment, error) {
	ctx := context.Background()

	_, err := db.NewInsert().Model(att).Exec(ctx)
	return *att, err
}

func UpdateAttachment(id int, att *Attachment) (Attachment, error) {
	ctx := context.Background()
	_, err := db.NewUpdate().Model(att).OmitZero().Where("id = ?", id).Exec(ctx)
	return *att, err
}

func GetAttachment(id int) (Attachment, error) {
	ctx := context.Background()
	table := new(Attachment)
	var att Attachment
	err := db.NewSelect().Model(table).Where("id = ?", id).Scan(ctx, &att)
	return att, err
}

func DeleteAttachment(id int) error {
	ctx := context.Background()
	table := new(Attachment)
	_, err := db.NewDelete().Model(table).Where("id = ?", id).Exec(ctx)

	return err
}
