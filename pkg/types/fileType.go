package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetadataModel struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	FileExtension string
	FileNameID    uuid.UUID `gorm:"type:uuid"`

	// * these are fields according which uuid has to be generated
	// todo think of something about object name type ex what if it contains russian letters
	ObjectName    string // * file name questions.txt
	MimeType      string // * file mime type text/plain, image/jpeg
	Size          int64
	StorageServer string
}

func (m *MetadataModel) BeforeCreate(tx *gorm.DB) error {
	// * we can customize uuid generation according to some fields
	m.ID = uuid.New()
	return nil
}
