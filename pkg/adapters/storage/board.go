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

func (r *boardRepo) GetUserBoards(ctx context.Context, userID uuid.UUID, limit, offset uint) (userBoards []board.Board, total uint, err error) {
	var int64Total int64
	var userBoardsEntities []entities.Board
	// Query to get the count of user boards
	userBoardsCountQuery := r.db.Table("boards").
		Joins("JOIN user_board_roles ubr ON ubr.board_id = boards.id").
		Where("ubr.user_id = ?", userID).
		Count(&int64Total)

	if userBoardsCountQuery.Error != nil {
		return nil, 0, userBoardsCountQuery.Error
	}

	// Query to get the boards where the user has a role
	userBoardsQuery := r.db.Table("boards").
		Select("boards.id, boards.name, boards.type, boards.created_at").
		Joins("JOIN user_board_roles ubr ON ubr.board_id = boards.id").
		Where("ubr.user_id = ?", userID).
		Order("boards.created_at DESC")

	if offset > 0 {
		userBoardsQuery = userBoardsQuery.Offset(int(offset))
	}

	if limit > 0 {
		userBoardsQuery = userBoardsQuery.Limit(int(limit))
	}

	if err := userBoardsQuery.Find(&userBoardsEntities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	total = uint(int64Total)
	userBoards = mappers.BatchBoardEntitiesToDomain(userBoardsEntities)
	return userBoards, total, nil
}

func (r *boardRepo) GetPublicBoards(ctx context.Context, userID uuid.UUID, limit, offset uint) (publicBoards []board.Board, total uint, err error) {
	var int64Total int64
	// Query to get the count of user boards
	publicBoardsCountQuery := r.db.Table("boards").
		Select("boards.id, boards.name, boards.type, boards.created_at").
		Where("boards.type = ? AND boards.id NOT IN (?)", "public",
			r.db.Table("user_board_roles").Select("board_id").Where("user_id = ?", userID)).
		Count(&int64Total)

	if publicBoardsCountQuery.Error != nil {
		return nil, 0, publicBoardsCountQuery.Error
	}

	// Query to get the public boards where the user does not have a role
	publicBoardsQuery := r.db.Table("boards").
		Select("boards.id, boards.name, boards.type, boards.created_at").
		Where("boards.type = ?", "public").
		Order("boards.created_at DESC")
	//Where("boards.type = ? AND boards.id NOT IN (?)", "public",
	//	r.db.Table("user_board_roles").Select("board_id").Where("user_id = ?", userID)).

	if offset > 0 {
		publicBoardsQuery = publicBoardsQuery.Offset(int(offset))
	}

	if limit > 0 {
		publicBoardsQuery = publicBoardsQuery.Limit(int(limit))
	}

	if err := publicBoardsQuery.Find(&publicBoards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	total = uint(int64Total)
	return publicBoards, total, nil
}

func (r *boardRepo) Insert(ctx context.Context, b *board.Board) error {
	boardEntity := mappers.BoardDomainToEntity(b)
	if err := r.db.WithContext(ctx).Save(&boardEntity).Error; err != nil {
		return err
	}

	b.ID = boardEntity.ID
	return nil
}

func (r *boardRepo) GetFullByID(ctx context.Context, id uuid.UUID) (*board.Board, error) {
	var b entities.Board

	if err := r.db.Preload("Users").
		Preload("Columns").
		Preload("Columns.Tasks").
		First(&b, "id = ?", id).Error; err != nil {
		return nil, err
	}
	domainBoard := mappers.BoardEntityToDomain(b)
	return &domainBoard, nil
}
