package storage

import (
	"context"
	"server/internal/column"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type columnRepo struct {
	db *gorm.DB
}

func NewColumnRepo(db *gorm.DB) column.Repo {
	return &columnRepo{
		db: db,
	}
}

func (r *columnRepo) Create(ctx context.Context, col *column.Column) (*column.Column, error) {
	if err := r.db.WithContext(ctx).Save(col).Error; err != nil {
		return nil, err
	}
	return col, nil
}

func (r *columnRepo) GetByID(ctx context.Context, id uuid.UUID) (*column.Column, error) {
	var col column.Column
	if err := r.db.WithContext(ctx).First(&col, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &col, nil
}

func (r *columnRepo) GetMaxOrderForBoard(ctx context.Context, boardID uuid.UUID) (uint, error) {
	var maxOrder uint
	err := r.db.WithContext(ctx).Model(&column.Column{}).Where("board_id = ?", boardID).Select("COALESCE(MAX(\"order\"), 0)").Scan(&maxOrder).Error
	return maxOrder, err
}

func (r *columnRepo) CreateBatch(ctx context.Context, columns []column.Column) ([]column.Column, error) {
	if err := r.db.WithContext(ctx).Create(&columns).Error; err != nil {
		return nil, err
	}
	return columns, nil
}

func (r *columnRepo) Delete(ctx context.Context, columnID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", columnID).Delete(&column.Column{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) GetByBoardID(ctx context.Context, boardID uuid.UUID) ([]column.Column, error) {
	var cols []column.Column
	if err := r.db.WithContext(ctx).Where("board_id = ?", boardID).Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}
