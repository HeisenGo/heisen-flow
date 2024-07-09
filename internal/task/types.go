/*
A self-referential relationship for subtasks is used:

ParentID and Parent for the parent task (null for top-level tasks).
Subtasks for child tasks.
*/

package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCircularDependency = errors.New("circular dependency detected")
)

type Repo interface {
	//GetUserTasks(ctx context.Context, userID uuid.UUID, limit, offset uint) ([]Board, uint, error)
	Insert(ctx context.Context, task *Task) error
	//GetByID(ctx context.Context, id uuid.UUID) (*Board, error)
}

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	// Status      TaskStatus `gorm:"not null"`
	Order           uint // in column which order is this
	StartAt         time.Time
	EndedAt         time.Time
	StoryPoint      uint      //(should be less than 10???)
	UserBoardRoleID uuid.UUID //Assignee
	CreatedByUBRID  uuid.UUID
	ColumnID        uuid.UUID
	BoardID         uuid.UUID
	ParentID        *uuid.UUID //can be null for tasks not sub tasks
	Subtasks        []Task

	DependsOn   []Task
	DependentBy []Task
}

type TaskDependency struct {
	DependentTaskID  uuid.UUID
	DependencyTaskID uuid.UUID
}
