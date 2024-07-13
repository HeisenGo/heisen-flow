package service

import (
	"context"
	"errors"
	"server/internal/board"
	"server/internal/column"
	userboardrole "server/internal/user_board_role"
	"server/pkg/adapters/storage/entities"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

var (
	ErrPermissionDeniedToCreateColumn = errors.New("permission denied: can not create column")
	ErrPermissionDeniedToDeleteColumn = errors.New("permission denied: can not delete the column")
)

type ColumnService struct {
	colOps           *column.Ops
	userBoardRoleOps *userboardrole.Ops
	boardOps         *board.Ops
}

func NewColumnService(colOps *column.Ops, userBoardRoleOps *userboardrole.Ops, boardOps *board.Ops) *ColumnService {
	return &ColumnService{colOps: colOps,
		boardOps:         boardOps,
		userBoardRoleOps: userBoardRoleOps}
}

func (s *ColumnService) CreateColumn(ctx context.Context, name string, boardID, userID uuid.UUID, order uint) (*entities.Column, error) {
	col := &column.Column{
		ID:       uuid.New(),
		Name:     name,
		BoardID:  boardID,
		OrderNum: order,
	}

	//check to see board exists?
	// check permission
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, col.BoardID)
	if err != nil {
		return nil, ErrPermissionDeniedToCreateColumn
	}

	if !rbac.HasPermission(role, rbac.PermissionManageColumns) {
		return nil, ErrPermissionDeniedToCreateColumn
	}

	if err := s.colOps.Create(ctx, col); err != nil {
		return nil, err
	}

	return &entities.Column{
		ID:       col.ID,
		Name:     col.Name,
		BoardID:  col.BoardID,
		OrderNum: col.OrderNum,
	}, nil
}

func (s *ColumnService) GetColumnByID(ctx context.Context, id uuid.UUID) (*entities.Column, error) {
	col, err := s.colOps.GetColumnByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entities.Column{
		ID:       col.ID,
		Name:     col.Name,
		BoardID:  col.BoardID,
		OrderNum: col.OrderNum,
	}, nil
}

func (s *ColumnService) GetMaxOrderForBoard(ctx context.Context, boardID uuid.UUID) (uint, error) {
	return s.colOps.GetMaxOrderForBoard(ctx, boardID)
}
func (s *ColumnService) GetMinOrderColumn(ctx context.Context, boardID uuid.UUID) (*entities.Column, error) {
	c, err := s.colOps.GetMinOrderColumn(ctx, boardID)
	if err != nil {
		return nil, err
	}
	return &entities.Column{
		ID:       c.ID,
		Name:     c.Name,
		BoardID:  c.BoardID,
		OrderNum: c.OrderNum,
	}, nil
}

func (s *ColumnService) CreateColumns(ctx context.Context, columns []entities.Column, userID uuid.UUID) ([]entities.Column, error) {
	colModels := make([]column.Column, len(columns))
	for i, col := range columns {
		colModels[i] = column.Column{
			ID:       col.ID,
			Name:     col.Name,
			BoardID:  col.BoardID,
			OrderNum: col.OrderNum,
		}
	}
	//check to see board exists?
	b, err := s.boardOps.GetBoardByID(ctx, colModels[0].BoardID)

	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, board.ErrBoardNotFound
	}
	// check permission
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, colModels[0].BoardID)
	if err != nil {
		return nil, ErrPermissionDeniedToCreateColumn
	}

	if !rbac.HasPermission(role, rbac.PermissionManageColumns) {
		return nil, ErrPermissionDeniedToCreateColumn
	}

	createdCols, err := s.colOps.CreateColumns(ctx, colModels)
	if err != nil {
		return nil, err
	}

	createdEntities := make([]entities.Column, len(createdCols))
	for i, col := range createdCols {
		createdEntities[i] = entities.Column{
			ID:       col.ID,
			Name:     col.Name,
			BoardID:  col.BoardID,
			OrderNum: col.OrderNum,
		}
	}
	return createdEntities, nil
}

func (s *ColumnService) DeleteColumn(ctx context.Context, columnID, userID uuid.UUID) error {

	// get column if exists
	col, err := s.colOps.GetColumnByID(ctx, columnID)
	if err != nil {
		return err
	}
	if col == nil {
		return column.ErrColumnNotFound
	}

	// check permission
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, col.BoardID)
	if err != nil {
		return ErrPermissionDeniedToDelete
	}

	if !rbac.HasPermission(role, rbac.PermissionManageColumns) {
		return ErrPermissionDeniedToDelete
	}
	return s.colOps.Delete(ctx, columnID)
}

func (s *ColumnService) ReorderColumns(ctx context.Context, userID, boardID uuid.UUID, newOrder map[uuid.UUID]uint) ([]column.Column, error) {

	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, boardID)
	if err != nil {
		return nil, ErrPermissionDenied
	}

	if !rbac.HasPermission(role, rbac.PermissionManageColumns) {
		return nil, ErrPermissionDeniedToDelete
	}
	err = s.colOps.ReorderColumns(ctx, boardID, newOrder)
	if err != nil {
		return nil, err
	}

	return s.colOps.GetColumns(ctx, boardID)
}
