package storage

import (
	"context"
	"server/internal/task"
	"server/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) task.Repo {
	return &taskRepo{
		db: db,
	}
}

func (r *taskRepo) Insert(ctx context.Context, t *task.Task) error {
	// check dependency circle [does the ids exist]
	// 
	// check roles 
	taskEntity := mappers.TaskDomainToEntity(t)
	if err := r.db.WithContext(ctx).Save(&taskEntity).Error; err != nil {
		return err
	}

	t.ID = taskEntity.ID
	return nil
}
