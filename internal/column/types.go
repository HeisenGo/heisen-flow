package column

import (
	"context"
	"errors"
	"regexp"
	"server/internal/task"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName          = errors.New("invalid column name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
	ErrColumnNotEmpty       = errors.New("you can't delete a column that has some tasks")
	ErrColumnNotFound       = errors.New("column doesn't exists")
	ErrFailedToFetchColumns = errors.New("failed to fetch columns")
	ErrInvalidColumnID      = errors.New("errInvalidColumnID")
	ErrFailedToUpdateColumn = errors.New("failed to update column")
	ErrLengthMismatch       = errors.New("length mismatch")
)

const (
	DoneDefaultColumn = "done"
)

type Repo interface {
	Create(ctx context.Context, column *Column) (*Column, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Column, error)
	GetMaxOrderForBoard(ctx context.Context, boardID uuid.UUID) (uint, error)
	GetMinOrderColumn(ctx context.Context, boardID uuid.UUID) (*Column, error)
	CreateBatch(ctx context.Context, columns []Column) ([]Column, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByBoardID(ctx context.Context, boardID uuid.UUID) ([]Column, error)
	SetDoneAsDefault(ctx context.Context, column *Column) error
	ReorderColumns(ctx context.Context, boardID uuid.UUID, newOrder map[uuid.UUID]uint) error
}

type Column struct {
	ID        uuid.UUID
	Name      string
	BoardID   uuid.UUID
	OrderNum  uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Tasks     []task.Task
}

func NewColumn(name string, boardID uuid.UUID, orderNum uint, createdAt time.Time) *Column {
	return &Column{
		Name:      name,
		BoardID:   boardID,
		OrderNum:  orderNum,
		CreatedAt: createdAt,
	}
}

func ValidateColumnName(name string) error {
	var validBoardName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validBoardName.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
