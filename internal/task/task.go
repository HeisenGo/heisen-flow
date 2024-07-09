/*
A self-referential relationship for subtasks is used:

ParentID and Parent for the parent task (null for top-level tasks).
Subtasks for child tasks.
*/

package entities

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID    uuid.UUID
	Title string
	// Status      TaskStatus `gorm:"not null"`
	Order           uint // in column which order is this
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
	EndedAt         time.Time
	StoryPoint      uint      //(should be less than 10???)
	UserBoardRoleID uuid.UUID //Assignee
	ColumnD         uuid.UUID
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
