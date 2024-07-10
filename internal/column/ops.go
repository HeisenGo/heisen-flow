package column

import (
	"context"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) GetColumnByID(ctx context.Context, id uuid.UUID) (*Column, error) {
	return o.repo.GetByID(ctx, id)
}

func (o *Ops) Create(ctx context.Context, column *Column) error {
	if err := ValidateColumnName(column.Name); err != nil {
		return err
	}
	_, err := o.repo.Create(ctx, column)
	return err
}

func (o *Ops) GetMaxOrderForBoard(ctx context.Context, boardID uuid.UUID) (uint, error) {
	return o.repo.GetMaxOrderForBoard(ctx, boardID)
}

func (o *Ops) CreateColumns(ctx context.Context, columns []Column) ([]Column, error) {
	for _, column := range columns {
		if err := ValidateColumnName(column.Name); err != nil {
			return nil, err
		}
	}
	return o.repo.CreateBatch(ctx, columns)
}
