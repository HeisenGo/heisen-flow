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
	StoryPoint  uint 

	// Relationships
	UserBoardRoleID *uuid.UUID    `gorm:"type:uuid"` //Assignee
	UserBoardRole   UserBoardRole `gorm:"foreignKey:UserBoardRoleID"`

	ColumnID uuid.UUID `gorm:"type:uuid"`
	Column   *Column   `gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE"`

	BoardID uuid.UUID `gorm:"type:uuid;not null"`
	Board   *Board    `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`

	ParentID *uuid.UUID `gorm:"type:uuid"` //can be null for tasks not sub tasks
	Parent   *Task      `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
	Subtasks []Task     `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`

	DependsOn   []Task `gorm:"many2many:task_dependencies;joinForeignKey:dependent_task_id;joinReferences:dependency_task_id;constraint:OnDelete:CASCADE"`
	DependentBy []Task `gorm:"many2many:task_dependencies;joinForeignKey:dependent_task_id;joinReferences:dependency_task_id;constraint:OnDelete:CASCADE"`
}

type TaskDependency struct {
	DependentTaskID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	DependencyTaskID uuid.UUID `gorm:"type:uuid;primaryKey"`
}
