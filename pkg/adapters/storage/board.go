package storage

import (
	"context"
	"errors"
	"server/internal/board"
	"server/pkg/adapters/storage/entities"
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

func (r *boardRepo) GetByID(ctx context.Context, id uuid.UUID) (*board.Board, error) {
	var b entities.Board

	err := r.db.WithContext(ctx).Model(&entities.Board{}).Where("id = ?", id).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	board := mappers.BoardEntityToDomain(b)
	return &board, nil
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
