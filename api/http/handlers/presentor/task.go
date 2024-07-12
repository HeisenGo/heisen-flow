package presenter

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
	"server/internal/task"
	"server/pkg/fp"
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

type CreateTaskResp struct {
	ID             uuid.UUID  `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	StartAt        time.Time  `json:"start_at"`
	EndAt          time.Time  `json:"end_at"`
	StoryPoint     uint       `json:"story_at"`
	AssigneeUserID *uuid.UUID `json:"assigneeUser_id"`
	ColumnID       uuid.UUID  `json:"column_id"`
	BoardID        uuid.UUID  `json:"board_id"`

	ParentID *uuid.UUID `json:"parentID"` //can be null for tasks not sub tasks

	DependsOn []DependTaskResp
}

type DependTaskResp struct {
	ID uuid.UUID `json:"id"`
}

func DomainTaskToDependTaskResp(task task.Task) DependTaskResp {
	return DependTaskResp{
		ID: task.ID,
	}
}

func BatchDomainTaskToDependTaskResp(tasks []task.Task) []DependTaskResp {
	return fp.Map(tasks, DomainTaskToDependTaskResp)
}

func DomainTaskToCreateTaskResp(task *task.Task) *CreateTaskResp {
	dependsOnTasks := BatchDomainTaskToDependTaskResp(task.DependsOn)
	return &CreateTaskResp{
		ID:             task.ID,
		Title:          task.Title,
		Description:    task.Description,
		StartAt:        task.StartAt,
		EndAt:          task.EndAt,
		StoryPoint:     task.StoryPoint,
		AssigneeUserID: task.AssigneeUserID,
		ColumnID:       task.ColumnID,
		BoardID:        task.BoardID,
		ParentID:       task.ParentID,
		DependsOn:      dependsOnTasks,
	}
}

func TaskToCreateTaskResp(task *task.Task) fiber.Map {
	return StructToFiberMap(DomainTaskToCreateTaskResp(task))
}

// Function to convert struct to fiber.Map
func StructToFiberMap(s interface{}) fiber.Map {
	result := fiber.Map{}
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := typ.Field(i)
		jsonTag := typeField.Tag.Get("json")

		// Skip if there's no json tag or it's set to "-"
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Set the field in the map
		result[jsonTag] = field.Interface()
	}

	return result
}
