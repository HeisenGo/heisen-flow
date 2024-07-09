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
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string    `gorm:"not null"`
	Description string
	// Status      TaskStatus `gorm:"not null"`
	Order      uint // in column which order is this
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	EndedAt    time.Time
	StoryPoint uint //(should be less than 10???)

	// Relationships
	UserID  uuid.UUID `gorm:"type:uuid"` //Assignee
	User    *User     `gorm:"foreignKey:UserID"`
	ColumnD uuid.UUID `gorm:"type:uuid"`
	//Column      *Column  !!!!!!!!!! need TO Do
	BoardID uuid.UUID `gorm:"type:uuid;not null"`
	Board   *Board    `gorm:"foreignKey:BoardID"`

	ParentID *uuid.UUID `gorm:"type:uuid"`
	Parent   *Task      `gorm:"foreignKey:ParentID"`
	Subtasks []Task     `gorm:"foreignKey:ParentID"`

	DependsOn   []Task `gorm:"many2many:task_dependencies;"`
	DependentBy []Task `gorm:"many2many:task_dependencies;joinForeignKey:dependent_task_id;joinReferences:dependency_task_id"`
}

// type TaskStatus string

// const (
// 	TaskStatusToDo       TaskStatus = "todo"
// 	TaskStatusInProgress TaskStatus = "in_progress"
// 	TaskStatusDone       TaskStatus = "done"
// )

type TaskDependency struct {
	DependentTaskID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	DependencyTaskID uuid.UUID `gorm:"type:uuid;primaryKey"`
}
