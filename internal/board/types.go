package board

import (
	"context"
	"errors"
	"regexp"
	"server/internal/column"
	"server/internal/user"
	"time"

	"github.com/google/uuid"
)

type BoardType string

const (
	Private BoardType = "private"
	Public  BoardType = "public"
)

var (
	ErrWrongType                      = errors.New("wrong type for board")
	ErrInvalidName                    = errors.New("invalid board name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
	ErrWrongBoardTime                 = errors.New("wrong board time")
	ErrBoardNotFound                  = errors.New("board not found")
	ErrFailedToDeleteBoard            = errors.New("failed to delete board")
	ErrFailedToFetchTasks             = errors.New("failed to fetch all tasks")
	ErrFailedToDeleteTaskDependencies = errors.New("failed to delete dependencies")
)

type Repo interface {
	Insert(ctx context.Context, board *Board) error
	GetByID(ctx context.Context, id uuid.UUID) (*Board, error)
	GetFullByID(ctx context.Context, id uuid.UUID) (*Board, error)
	GetUserBoards(ctx context.Context, userID uuid.UUID, limit, offset uint) (userBoards []Board, total uint, err error)
	GetPublicBoards(ctx context.Context, userID uuid.UUID, limit, offset uint) (publicBoards []Board, total uint, err error)
	DeleteByID(ctx context.Context, boardID uuid.UUID) error
}

type Board struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Name      string
	Type      string
	Users     []user.User
	Columns   []column.Column
}

func ValidateBoardName(name string) error {
	var validBoardName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validBoardName.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
