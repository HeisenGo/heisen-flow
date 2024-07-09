package storage

import (
	"context"
	"server/internal/board"
	"server/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
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

func (r *boardRepo) GetUserBoards(ctx context.Context, userID uuid.UUID, limit, offset uint) ([]board.Board, uint, error) {
	return nil, 0, nil
}

func (r *boardRepo) Insert(ctx context.Context, b *board.Board) error {
	boardEntity := mappers.BoardDomainToEntity(b)
	if err := r.db.WithContext(ctx).Save(&boardEntity).Error; err != nil {
		return err
	}

	b.ID = boardEntity.ID
	return nil
}
