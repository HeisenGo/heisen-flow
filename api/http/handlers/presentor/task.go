package presenter

import (
	"server/internal/task"
	"time"

	"github.com/google/uuid"
)

type UserTask struct {
	ID             uuid.UUID `json:"task_id"`
	BoardID        uuid.UUID `json:"board_id" validate:"required"`
	StartAt        Timestamp `json:"start_at"`
	EndAt          Timestamp `json:"end_at"`
	AssigneeUserID uuid.UUID `json:"assignee_user_id" validate:"required"`
	Title          string    `json:"title" validate:"required"`
	Description    string    `json:"desc"`
	StoryPoint     uint      `json:"story_point"`
	// for tasks that this task depends on
	DependsOnTaskIDs []uuid.UUID `json:"depends_on_task_ids"`
	//for tasks that depend on this task
	ParentID *uuid.UUID `json:"parent_id"`
}

type DependentTasks struct {
	ID               uuid.UUID   `json:"task_id" validate:"required"`
	DependsOnTaskIDs []uuid.UUID `json:"depends_on_task_ids" validate:"required"`
}

func AddDependencyReqToTask(dependentTasksReq *DependentTasks, userID uuid.UUID) *task.Task {
	return &task.Task{
		ID:               dependentTasksReq.ID,
		DependsOnTaskIDs: dependentTasksReq.DependsOnTaskIDs,
		CreatedByUserID:  userID,
	}
}

func TaskToUserTask(t task.Task) UserTask {
	return UserTask{
		ID:          t.ID,
		Description: t.Description,
	}
}

func UserTaskToTask(userTaskReq *UserTask, userID uuid.UUID) *task.Task {
	return &task.Task{
		Title:            userTaskReq.Title,
		Description:      userTaskReq.Description,
		StartAt:          time.Time(userTaskReq.StartAt),
		EndAt:            time.Time(userTaskReq.EndAt),
		StoryPoint:       userTaskReq.StoryPoint,
		BoardID:          userTaskReq.BoardID,
		CreatedByUserID:  userID,
		ParentID:         userTaskReq.ParentID,
		DependsOnTaskIDs: userTaskReq.DependsOnTaskIDs,
	}
}
