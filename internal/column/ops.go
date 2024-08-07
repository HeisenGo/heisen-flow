package column

import (
	"context"

	"time"

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
func (o *Ops) GetMinOrderColumn(ctx context.Context, boardID uuid.UUID) (*Column, error) {
	return o.repo.GetMinOrderColumn(ctx, boardID)
}

func (o *Ops) CreateColumns(ctx context.Context, columns []Column) ([]Column, error) {
	for _, column := range columns {
		if err := ValidateColumnName(column.Name); err != nil {
			return nil, err
		}
	}
	return o.repo.CreateBatch(ctx, columns)
}

func (o *Ops) Delete(ctx context.Context, columnID uuid.UUID) error {
	return o.repo.Delete(ctx, columnID)
}

func (o *Ops) GetColumnsByBoardID(ctx context.Context, boardID uuid.UUID) ([]Column, error) {
	return o.repo.GetByBoardID(ctx, boardID)
}

func (o *Ops) SetDoneAsDefault(ctx context.Context, boardID uuid.UUID) (*Column, error) {
	col := NewColumn(DoneDefaultColumn, boardID, uint(1), time.Now())
	err := o.repo.SetDoneAsDefault(ctx, col)
	if err != nil {
		return nil, err
	}
	return col, nil
}

func (o *Ops) ReorderColumns(ctx context.Context, boardID uuid.UUID, newOrder map[uuid.UUID]uint) error {
	return o.repo.ReorderColumns(ctx, boardID, newOrder)
}

func (o *Ops) GetColumns(ctx context.Context, boardID uuid.UUID) ([]Column, error) {
	return o.repo.GetColumns(ctx, boardID)
}
