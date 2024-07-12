package mappers

import (
	"server/internal/column"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"

	"gorm.io/gorm"
)

func ColumnEntityToDomain(col entities.Column) column.Column {
	tasks := BatchTaskEntitiesToDomain(col.Tasks)
	return column.Column{
		ID:        col.ID,
		Name:      col.Name,
		BoardID:   col.BoardID,
		OrderNum:  col.OrderNum,
		CreatedAt: col.CreatedAt,
		UpdatedAt: col.UpdatedAt,
		Tasks:     tasks,
	}
}

func BatchColumnEntitiesToDomain(cols []entities.Column) []column.Column {
	return fp.Map(cols, ColumnEntityToDomain)
}

func ColumnDomainToEntity(col column.Column) entities.Column {
	return entities.Column{
		ID:        col.ID,
		Name:      col.Name,
		BoardID:   col.BoardID,
		OrderNum:  col.OrderNum,
		CreatedAt: col.CreatedAt,
		UpdatedAt: col.UpdatedAt,
		DeletedAt: gorm.DeletedAt{},
	}
}

func ColumnDomainsToEntities(cols []column.Column) []entities.Column {
	return fp.Map(cols, ColumnDomainToEntity)
}
