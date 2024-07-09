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
	if err := r.db.WithContext(ctx).Create(col).Error; err != nil {
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
