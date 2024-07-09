package service

import (
	"context"
	"server/internal/column"

	"github.com/google/uuid"
)

type ColumnService struct {
	colOps *column.Ops
}

func NewColumnService(colOps *column.Ops) *ColumnService {
	return &ColumnService{colOps: colOps}
}

func (s *ColumnService) CreateColumn(ctx context.Context, name string, boardID uuid.UUID, order uint) (*column.Column, error) {
	col := &column.Column{
		ID:      uuid.New(),
		Name:    name,
		BoardID: boardID,
		Order:   order,
	}

	if err := s.colOps.Create(ctx, col); err != nil {
		return nil, err
	}

	return col, nil
}

func (s *ColumnService) GetColumnByID(ctx context.Context, id uuid.UUID) (*column.Column, error) {
	col, err := s.colOps.GetColumnByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return col, nil
}
