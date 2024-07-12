/*
A self-referential relationship for subtasks is used:

ParentID and Parent for the parent task (null for top-level tasks).
Subtasks for child tasks.
*/

package task

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCircularDependency             = errors.New("circular dependency detected")
	ErrFailedToFindDependsOnTasks     = errors.New("failed to find depends on tasks")
	ErrFailedToCreateTaskDependencies = errors.New("failed to create task dependencies")
	ErrEmptyTitle                     = errors.New("title is required")
	ErrLongTitle                      = errors.New("title cannot be longer than 255 characters")
	ErrLongDescription                = errors.New("description cannot be longer than 1000 characters")
	ErrTitleInvalidCharacter          = errors.New("title contains invalid characters")
	ErrDescInvalidCharacter           = errors.New("description contains invalid characters")
	ErrParentTaskNotFound             = errors.New("parent not found")
	ErrTaskNotFound                   = errors.New("task not found")
	ErrBoardNotFound                  = errors.New("board not found")
	ErrInvalidStoryPoint              = errors.New("invalid story point: must be one of 1, 2, 3, 5, 8, 13, 21")
)

type Repo interface {
	Insert(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	AddDependency(ctx context.Context, t *Task) error
}

type Task struct {
	ID              uuid.UUID
	Title           string
	Description     string
	Order           uint // in column which order is this
	StartAt         time.Time
	EndAt           time.Time
	StoryPoint      uint
	AssigneeUserID  *uuid.UUID
	UserBoardRoleID *uuid.UUID //Assignee
	CreatedByUserID uuid.UUID
	ColumnID        uuid.UUID
	BoardID         uuid.UUID

	ParentID   *uuid.UUID //can be null for tasks not sub tasks
	SubTaskIDs []uuid.UUID
	Subtasks   []Task

	DependsOn          []Task
	DependsOnTaskIDs   []uuid.UUID
	DependentBy        []Task
	DependentByTaskIDs []uuid.UUID
}

type TaskDependency struct {
	DependentTaskID  uuid.UUID
	DependencyTaskID uuid.UUID
}

func validateTitleAndDescription(title, description string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if len(title) > 255 {
		return ErrLongTitle
	}
	if len(description) > 3000 {
		return ErrLongDescription
	}

	invalidChars := []string{";", "--", "'"}
	for _, char := range invalidChars {
		if strings.Contains(title, char) {
			return ErrTitleInvalidCharacter
		}
		if strings.Contains(description, char) {
			return ErrDescInvalidCharacter
		}
	}

	return nil
}

func validateStoryPoint(storyPoint uint) error {
	allowedStoryPoints := []uint{1, 2, 3, 5, 8, 13, 21}
	for _, sp := range allowedStoryPoints {
		if storyPoint == sp {
			return nil
		}
	}
	return ErrInvalidStoryPoint
}
