package column

import (
	"context"
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var (
	ErrInvalidName = errors.New("invalid column name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
)

type Repo interface {
	Create(ctx context.Context, column *Column) (*Column, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Column, error)
	GetMaxOrderForBoard(ctx context.Context, boardID uuid.UUID) (uint, error)
	CreateBatch(ctx context.Context, columns []Column) ([]Column, error)
}

type Column struct {
	ID      uuid.UUID
	Name    string
	BoardID uuid.UUID
	Order   uint
}

func ValidateColumnName(name string) error {
	var validBoardName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validBoardName.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
