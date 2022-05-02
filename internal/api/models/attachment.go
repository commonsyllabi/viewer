package models

import "github.com/commonsyllabi/viewer/pkg/commoncartridge"

type Attachment struct {
	id        int64 `bun:",pk,autoincrement"`
	name      string
	cartridge commoncartridge.Cartridge
}
