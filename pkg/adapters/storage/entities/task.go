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
	Order       uint // in column which order is this
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	StartAt     time.Time
	EndAt       time.Time
	StoryPoint  uint //(should be less than 1 2 3 5 8 13 21 ???)

	// Relationships
	UserBoardRoleID *uuid.UUID     `gorm:"type:uuid"` //Assignee
	UserBoardRole   *UserBoardRole `gorm:"foreignKey:UserBoardRoleID"`

	ColumnID uuid.UUID `gorm:"type:uuid"`
	//Column      *Column  !!!!!!!!!! need TO Do

	BoardID uuid.UUID `gorm:"type:uuid;not null"`
	Board   *Board    `gorm:"foreignKey:BoardID"`

	ParentID *uuid.UUID `gorm:"type:uuid"` //can be null for tasks not sub tasks
	Parent   *Task      `gorm:"foreignKey:ParentID"`
	Subtasks []Task     `gorm:"foreignKey:ParentID"`
	// for tasks that this task depends on
	DependsOn []*Task `gorm:"many2many:task_dependencies;joinForeignKey:dependent_task_id;joinReferences:dependency_task_id"`
	// for tasks that depend on this task.
	DependentBy []*Task `gorm:"many2many:task_dependencies;joinForeignKey:dependency_task_id;joinReferences:dependent_task_id"`
}

type TaskDependency struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	DependentTaskID  uuid.UUID      `gorm:"index:idx_together_dependent_dependency,unique"`
	DependencyTaskID uuid.UUID      `gorm:"index:idx_together_dependent_dependency,unique"`
}
