package service

import (
	"context"
	"server/internal/column"
	"server/pkg/adapters/storage/entities"

	"github.com/google/uuid"
)

type ColumnService struct {
	colOps *column.Ops
}

func NewColumnService(colOps *column.Ops) *ColumnService {
	return &ColumnService{colOps: colOps}
}

func (s *ColumnService) CreateColumn(ctx context.Context, name string, boardID uuid.UUID, order uint) (*entities.Column, error) {
	col := &column.Column{
		ID:      uuid.New(),
		Name:    name,
		BoardID: boardID,
		Order:   order,
	}

	if err := s.colOps.Create(ctx, col); err != nil {
		return nil, err
	}

	return &entities.Column{
		ID:      col.ID,
		Name:    col.Name,
		BoardID: col.BoardID,
		Order:   col.Order,
	}, nil
}

func (s *ColumnService) GetColumnByID(ctx context.Context, id uuid.UUID) (*entities.Column, error) {
	col, err := s.colOps.GetColumnByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entities.Column{
		ID:      col.ID,
		Name:    col.Name,
		BoardID: col.BoardID,
		Order:   col.Order,
	}, nil
}

func (s *ColumnService) GetMaxOrderForBoard(ctx context.Context, boardID uuid.UUID) (uint, error) {
	return s.colOps.GetMaxOrderForBoard(ctx, boardID)
}

func (s *ColumnService) CreateColumns(ctx context.Context, columns []entities.Column) ([]entities.Column, error) {
	colModels := make([]column.Column, len(columns))
	for i, col := range columns {
		colModels[i] = column.Column{
			ID:      col.ID,
			Name:    col.Name,
			BoardID: col.BoardID,
			Order:   col.Order,
		}
	}

	createdCols, err := s.colOps.CreateColumns(ctx, colModels)
	if err != nil {
		return nil, err
	}

	createdEntities := make([]entities.Column, len(createdCols))
	for i, col := range createdCols {
		createdEntities[i] = entities.Column{
			ID:      col.ID,
			Name:    col.Name,
			BoardID: col.BoardID,
			Order:   col.Order,
		}
	}
	return createdEntities, nil
}

func (s *ColumnService) DeleteColumn(ctx context.Context, columnID uuid.UUID) error {
	return s.colOps.Delete(ctx, columnID)
}

// func (s *ColumnService) GetColumnsByBoardID(ctx context.Context, boardID uuid.UUID) ([]column.Column, error) {
// 	return s.colOps.GetColumnsByBoardID(ctx, boardID)
// }
