package mappers

import (
	"server/internal/board"
	"server/internal/task"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"
)

func TaskEntityToDomain(boardEntity entities.Board) board.Board {
	return board.Board{
		ID:        boardEntity.ID,
		CreatedAt: boardEntity.CreatedAt,
		Name:      boardEntity.Name,
		Type:      boardEntity.Type,
	}
}

func TaskEntitiesToDomain(boardEntities []entities.Board) []board.Board {
	return fp.Map(boardEntities, BoardEntityToDomain)
}

func TaskDomainToEntity(t *task.Task) *entities.Task {
	return &entities.Task{
		Title: t.Title,
		Description: t.Description,
		StartAt: t.StartAt,
		EndAt: t.EndAt,
		StoryPoint: t.StoryPoint,
		UserBoardRoleID: t.UserBoardRoleID,
		BoardID: t.BoardID,
		ParentID: t.ParentID,
	}
}
