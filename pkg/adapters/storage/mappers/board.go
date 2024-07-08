package mappers

import (
	"server/internal/board"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"
)

func BoardEntityToDomain(boardEntity entities.Board) board.Board {
	return board.Board{
		ID:        boardEntity.ID,
		CreatedAt: boardEntity.CreatedAt,
		Name:      boardEntity.Name,
		Type:      boardEntity.Type,
	}
}

func BoardEntitiesToDomain(boardEntities []entities.Board) []board.Board {
	return fp.Map(boardEntities, BoardEntityToDomain)
}

func BoardDomainToEntity(b *board.Board) *entities.Board {
	return &entities.Board{
		CreatedAt: b.CreatedAt,
		Name:      b.Name,
		Type:      b.Type,
	}
}
