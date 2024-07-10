package mappers

import (
	"server/internal/board"
	"server/internal/task"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"
)

func TaskEntityToDomain(taskEntity entities.Task) task.Task {
	return task.Task{
		ID:              taskEntity.ID,
		Title:           taskEntity.Title,
		Description:     taskEntity.Description,
		StartAt:         taskEntity.StartAt,
		EndAt:           taskEntity.EndAt,
		StoryPoint:      taskEntity.StoryPoint,
		UserBoardRoleID: taskEntity.UserBoardRoleID,
		BoardID:         taskEntity.BoardID,
		ParentID:        taskEntity.ParentID,
	}
}

func TaskEntitiesToDomain(boardEntities []entities.Board) []board.Board {
	return fp.Map(boardEntities, BoardEntityToDomain)
}

func TaskDomainToEntity(t *task.Task) *entities.Task {
	return &entities.Task{
		Title:           t.Title,
		Description:     t.Description,
		StartAt:         t.StartAt,
		EndAt:           t.EndAt,
		StoryPoint:      t.StoryPoint,
		UserBoardRoleID: t.UserBoardRoleID,
		BoardID:         t.BoardID,
		ParentID:        t.ParentID,
	}
}
