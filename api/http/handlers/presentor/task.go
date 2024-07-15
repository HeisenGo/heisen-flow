package presenter

import (
	"server/internal/comment"
	"server/internal/task"
	"server/internal/user"
	"server/pkg/fp"
	"time"

	"github.com/google/uuid"
)

type UserTask struct {
	ID             uuid.UUID  `json:"task_id"`
	BoardID        uuid.UUID  `json:"board_id" validate:"required"`
	StartAt        *time.Time `json:"start_at"`
	EndAt          *time.Time `json:"end_at"`
	AssigneeUserID uuid.UUID  `json:"assignee_user_id" validate:"required"`
	Title          string     `json:"title" validate:"required"`
	Description    string     `json:"desc"`
	StoryPoint     uint       `json:"story_point"`
	// for tasks that this task depends on
	DependsOnTaskIDs []uuid.UUID `json:"depends_on_task_ids"`
	//for tasks that depend on this task
	ParentID *uuid.UUID `json:"parent_id"`
}

type ReorderTasksReq struct {
	ColumnID uuid.UUID         `json:"column_id"`
	Tasks    []ReorderTaskItem `json:"tasks"`
}

type ReorderTaskItem struct {
	ID uuid.UUID `json:"id"`
}

func ReorderTasksReqToMap(req ReorderTasksReq) (uuid.UUID, map[uuid.UUID]uint) {
	newOrder := make(map[uuid.UUID]uint)
	for i, t := range req.Tasks {
		newOrder[t.ID] = uint(i + 1)
	}
	return req.ColumnID, newOrder
}

type TaskReorderRespItem struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Order uint      `json:"order"`
}

func TaskToTaskReorderResp(t task.Task) TaskReorderRespItem {
	return TaskReorderRespItem{
		ID:    t.ID,
		Title: t.Title,
		Order: t.Order,
	}
}
func BatchTaskToTaskReorderRespItem(cols []task.Task) []TaskReorderRespItem {
	return fp.Map(cols, TaskToTaskReorderResp)
}

type DependentTasks struct {
	ID               uuid.UUID   `json:"task_id" validate:"required"`
	DependsOnTaskIDs []uuid.UUID `json:"depends_on_task_ids" validate:"required"`
}

type UpdateTaskColReq struct {
	ColumnID uuid.UUID `json:"column_id" validate:"required"`
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
		StartAt:          userTaskReq.StartAt,
		EndAt:            userTaskReq.EndAt,
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
type TaskCommentResp struct {
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
}

type FullTaskResp struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Order       uint       `json:"order"`
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	StoryPoint  uint       `json:"story_point"`

	// Relationships
	User     *TaskUserResp     `json:"user"`
	Parent   *TaskParentResp   `json:"parent"`
	Subtasks []TaskSubTaskResp `json:"subtasks"`
	//TODO:Comments []Comment  `gorm:"foreignKey:TaskID"`

	DependsOn []TaskDependTaskResp `json:"dependencies"`
	Comments  []TaskCommentResp    `json:"comments"`
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
func CommentToTaskCommentResp(c comment.Comment) TaskCommentResp {
	return TaskCommentResp{
		CreatedAt:   c.CreatedAt,
		Description: c.Description,
		Title:       c.Title,
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
func BatchCommentToTaskCommentResp(comments []comment.Comment) []TaskCommentResp {
	return fp.Map(comments, CommentToTaskCommentResp)
}

type UpdatedTaskResp struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Order       uint       `json:"order"`
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	StoryPoint  uint       `json:"story_point"`
}

func TaskToUpdatedTaskResp(t task.Task) UpdatedTaskResp {
	return UpdatedTaskResp{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Order:       t.Order,
		StartAt:     t.StartAt,
		EndAt:       t.EndAt,
		StoryPoint:  t.StoryPoint,
	}
}

func TaskToFullTaskResp(t task.Task) FullTaskResp {
	var (
		p          *TaskParentResp
		subs       []TaskSubTaskResp
		dependsOns []TaskDependTaskResp
		comments   []TaskCommentResp
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
	if t.Comments != nil {
		comments = BatchCommentToTaskCommentResp(t.Comments)

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
		Comments:    comments,
	}
}

type CreateTaskResp struct {
	ID             uuid.UUID  `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	StartAt        *time.Time `json:"start_at"`
	EndAt          *time.Time `json:"end_at"`
	StoryPoint     uint       `json:"story_at"`
	AssigneeUserID *uuid.UUID `json:"assignee_user_id"`
	ColumnID       uuid.UUID  `json:"column_id"`
	BoardID        uuid.UUID  `json:"board_id"`

	ParentID *uuid.UUID `json:"parent_id"` //can be null for tasks not sub tasks

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
