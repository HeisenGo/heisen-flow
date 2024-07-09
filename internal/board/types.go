package board

import (
	"context"
	"errors"
	"regexp"
	"time"
	"github.com/google/uuid"
)

type BoardType string

const (
	Private BoardType = "private"
	Public  BoardType = "public"
)

var (
	ErrWrongType      = errors.New("wrong type for board")
	ErrInvalidName    = errors.New("invalid board name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
	ErrWrongBoardTime = errors.New("wrong board time")
	ErrBoardNotFound = errors.New("board not found")
)

type Repo interface {
	GetUserBoards(ctx context.Context, userID uuid.UUID, limit, offset uint) ([]Board, uint, error)
	Insert(ctx context.Context, board *Board) error
	GetByID(ctx context.Context, id uuid.UUID) (*Board, error)}

type Board struct {
	// UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt
	ID        uuid.UUID
	CreatedAt time.Time
	Name      string
	Type      string
}

func ValidateBoardName(name string) error {
	var validBoardName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validBoardName.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
