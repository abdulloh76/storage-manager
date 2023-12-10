package types

import "github.com/google/uuid"

type MetadataStore interface {
	CreateMetadata(fileMetadata *MetadataModel) (id uuid.UUID, err error)
	GetMetadata(id uuid.UUID) (*MetadataModel, error)
}
