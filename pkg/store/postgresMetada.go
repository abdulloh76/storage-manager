package store

import (
	"github.com/abdulloh76/storage-manager/pkg/types"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresMetadataStore struct {
	db *gorm.DB
}

var _ types.MetadataStore = (*PostgresMetadataStore)(nil)

func NewPostgresDBStore(dbURL string) *PostgresMetadataStore {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&types.MetadataModel{})

	return &PostgresMetadataStore{
		db,
	}
}

func (m *PostgresMetadataStore) CreateMetadata(fileMetadata *types.MetadataModel) (id uuid.UUID, err error) {
	err = m.db.Create(&fileMetadata).Error

	return fileMetadata.ID, err
}

func (m *PostgresMetadataStore) GetMetadata(id uuid.UUID) ([]types.MetadataModel, error) {
	var metadata []types.MetadataModel
	err := m.db.Where("file_name_id = ?", id).Find(&metadata).Error

	return metadata, err
}
