package test

import (
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type MockUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type MockUserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MockBoard struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserCreationResult struct {
	StatusCode int
	Message    string
}

type MockTask struct {
	Title            string      `json:"title" validate:"required"`
	Description      string      `json:"desc"`
	StartAt          *time.Time  `json:"start_at"`
	EndAt            *time.Time  `json:"end_at"`
	AssigneeUserID   uuid.UUID   `json:"assignee_user_id" validate:"required"`
	StoryPoint       uint        `json:"story_point"`
	DependsOnTaskIDs []uuid.UUID `json:"depends_on_task_ids"`
	ParentID         *uuid.UUID  `json:"parent_id"`
	BoardID          uuid.UUID   `json:"board_id"`
}

type UserCreationData struct {
	UserID string `json:"user_id"`
}

type BoardCreationData struct {
	BoardID string `json:"board_id"`
}

type Column struct {
	ID      uuid.UUID
	Name    string
	BoardID uuid.UUID
}
