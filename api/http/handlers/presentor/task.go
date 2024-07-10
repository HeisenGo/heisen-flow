package presenter

import (
	"server/internal/task"
	"time"

	"github.com/google/uuid"
)

type UserTask struct {
	ID             uuid.UUID `json:"task_id"`
	StartAt        Timestamp `json:"start_at"`
	EndAt          Timestamp `json:"end_at"`
	AssigneeUserID *uuid.UUID `json:"assignee_user_id"`
	Title          string    `json:"title"`
	Description    string    `json:"desc"`
	StoryPoint     uint      `json:"story_point"`
	// for tasks that this task depends on
	DependsOnTaskIDs []uuid.UUID `json:"depends_on_task_ids"`
	//for tasks that depend on this task
	// DependentByTaskIDs []uuid.UUID `json:"dependent_by_task_ids"`
	ParentID *uuid.UUID `json:"parent_id"`
	// SubTasksIDs        []uuid.UUID `json:"sub_tasks"`
}

func TaskToUserTask(t task.Task) UserTask {
	return UserTask{
		ID: t.ID,
		// TotalPrice:    o.TotalPrice,
		// TotalQuantity: o.TotalQuantity,
		Description: t.Description,
	}
}

// func OrdersToUserOrders(orders []order.Order) []UserOrder {
// 	return fp.Map(orders, OrderToUserOrder)
// }

func UserTaskToTask(userTaskReq *UserTask, userID uuid.UUID) *task.Task {
	return &task.Task{
		Title:           userTaskReq.Title,
		Description:     userTaskReq.Description,
		StartAt:         time.Time(userTaskReq.StartAt),
		EndAt:           time.Time(userTaskReq.EndAt),
		StoryPoint:      userTaskReq.StoryPoint,
		CreatedByUserID: userID, //
		// BoardID: t.BoardID,
		ParentID:         userTaskReq.ParentID,
		DependsOnTaskIDs: userTaskReq.DependsOnTaskIDs,
		AssigneeUserID: userTaskReq.AssigneeUserID, // this will be changed during 
		// DependentByTaskIDs: userTaskReq.DependentByTaskIDs,
		// SubTaskIDs:         userTaskReq.SubTasksIDs,
	}
}
