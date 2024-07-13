package comment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	Insert(ctx context.Context, comment *Comment) error
}

type Comment struct {
	ID              uuid.UUID
	Title           string
	Description     string
	UserBoardRoleID uuid.UUID
	TaskID          uuid.UUID
	CreatedAt       time.Time
}