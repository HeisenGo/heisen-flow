package storage

import (
	"context"
	"server/internal/board"
	"server/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type boardRepo struct {
	db *gorm.DB
}

func NewBoardRepo(db *gorm.DB) board.Repo {
	return &boardRepo{
		db: db,
	}
}

func (r *boardRepo) GetUserBoards(ctx context.Context, userID uint, limit, offset uint) ([]board.Board, uint, error) {
	return nil, 0, nil
}

func (r *boardRepo) Insert(ctx context.Context, o *board.Board) error {
	boardEntity := mappers.BoardDomainToEntity(o)
	if err := r.db.WithContext(ctx).Save(&boardEntity).Error; err != nil {
		return err
	}

	o.ID = boardEntity.ID
	return nil
}
