package mappers

import (
	"server/internal/column"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"

	"gorm.io/gorm"
)

func ColumnEntityToDomain(col entities.Column) column.Column {
	return column.Column{
		ID:        col.ID,
		Name:      col.Name,
		BoardID:   col.BoardID,
		Order:     col.Order,
		CreatedAt: col.CreatedAt,
		UpdatedAt: col.UpdatedAt,
	}
}

func ColumnEntitiesToDomain(cols []entities.Column) []column.Column {
	return fp.Map(cols, ColumnEntityToDomain)
}

func ColumnDomainToEntity(col column.Column) entities.Column {
	return entities.Column{
		ID:        col.ID,
		Name:      col.Name,
		BoardID:   col.BoardID,
		Order:     col.Order,
		CreatedAt: col.CreatedAt,
		UpdatedAt: col.UpdatedAt,
		DeletedAt: gorm.DeletedAt{},
	}
}

func ColumnDomainsToEntities(cols []column.Column) []entities.Column {
	return fp.Map(cols, ColumnDomainToEntity)
}
