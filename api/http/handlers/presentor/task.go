package presenter

import (
	"reflect"
	"server/internal/task"
	"server/internal/user"
	"server/pkg/fp"
	"time"

	"github.com/gofiber/fiber/v2"
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
		AssigneeUserID:   &userTaskReq.AssigneeUserID,
	}
}

type TaskUserResp struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type TaskColumnResp struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TaskBoardResp struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TaskParentResp struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type TaskSubTaskResp struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type TaskDependTaskResp struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type FullTaskResp struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Order       uint      `json:"order"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	StoryPoint  uint      `json:"story_point"`

	// Relationships
	User     *TaskUserResp     `json:"user"`
	Parent   *TaskParentResp   `json:"parent"`
	Subtasks []TaskSubTaskResp `json:"subtasks"`
	//TODO:Comments []Comment  `gorm:"foreignKey:TaskID"`

	DependsOn []TaskDependTaskResp `json:"depends_on"`
}

func UserToTaskUserResp(u user.User) TaskUserResp {
	return TaskUserResp{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

func TaskToTaskSubTaskResp(t task.Task) TaskSubTaskResp {
	return TaskSubTaskResp{
		ID:    t.ID,
		Title: t.Title,
	}
}

func TaskToTaskDependTaskResp(t task.Task) TaskDependTaskResp {
	return TaskDependTaskResp{
		ID:    t.ID,
		Title: t.Title,
	}
}

func TaskToTaskParentResp(t task.Task) *TaskParentResp {
	return &TaskParentResp{
		ID:    t.ID,
		Title: t.Title,
	}
}

func BatchTaskToTaskSubTaskResp(tasks []task.Task) []TaskSubTaskResp {
	return fp.Map(tasks, TaskToTaskSubTaskResp)
}

func BatchTaskToTaskDependTaskResp(tasks []task.Task) []TaskDependTaskResp {
	return fp.Map(tasks, TaskToTaskDependTaskResp)
}

func TaskToFullTaskResp(t task.Task) FullTaskResp {
	var (
		p          *TaskParentResp
		subs       []TaskSubTaskResp
		dependsOns []TaskDependTaskResp
	)

	u := UserToTaskUserResp(*t.UserBoardRole.User)
	if t.Parent != nil {
		p = TaskToTaskParentResp(*t.Parent)
	}
	if t.Subtasks != nil {
		subs = BatchTaskToTaskSubTaskResp(t.Subtasks)
	}
	if t.DependsOn != nil {
		dependsOns = BatchTaskToTaskDependTaskResp(t.DependsOn)

	}
	return FullTaskResp{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Order:       t.Order,
		StartAt:     t.StartAt,
		EndAt:       t.EndAt,
		StoryPoint:  t.StoryPoint,
		User:        &u,
		Parent:      p,
		Subtasks:    subs,
		DependsOn:   dependsOns,
	}
}

type CreateTaskResp struct {
	ID             uuid.UUID  `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	StartAt        time.Time  `json:"start_at"`
	EndAt          time.Time  `json:"end_at"`
	StoryPoint     uint       `json:"story_at"`
	AssigneeUserID *uuid.UUID `json:"assignee_user_id"`
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
