package mappers

import (
	"server/internal/task"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"

	"github.com/google/uuid"
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

func BatchTaskEntitiesToDomain(taskEntities []entities.Task) []task.Task {
	return fp.Map(taskEntities, TaskEntityToDomain)
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

func TaskDependencyDomainToTaskEntity(id uuid.UUID) entities.Task {
	return entities.Task{ID: id}
}

func BatchTaskDependencyDomainToTask(taskIDs []uuid.UUID) []entities.Task {
	tasks := make([]entities.Task, len(taskIDs))

	for _, id := range taskIDs {
		t := TaskDependencyDomainToTaskEntity(id)
		tasks = append(tasks, t)
	}
	return tasks
}