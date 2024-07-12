package storage

import (
	"context"
	"gorm.io/gorm"
	"server/internal/comment"
	"server/pkg/adapters/storage/mappers"
)

type commentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) comment.Repo {
	return &commentRepo{
		db: db,
	}
}

func (r *commentRepo) Insert(ctx context.Context, comment *comment.Comment) error {
	commentEntity := mappers.CommentDomainToEntity(comment)
	if err := r.db.WithContext(ctx).Save(&commentEntity).Error; err != nil {
		return err
	}

	comment.ID = comment.ID
	return nil
}
