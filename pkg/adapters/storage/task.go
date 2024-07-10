package storage

import (
	"context"
	"errors"
	"fmt"
	"server/internal/task"
	"server/pkg/adapters/storage/entities"
	"server/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
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

func (r *taskRepo) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	var t entities.Task

	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", id).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	task := mappers.TaskEntityToDomain(t)
	return &task, nil
}

func (r *taskRepo) Insert(ctx context.Context, t *task.Task) error {
	taskEntity := mappers.TaskDomainToEntity(t)

	if err := r.db.WithContext(ctx).Save(&taskEntity).Error; err != nil {
		return err
	}

	t.ID = taskEntity.ID
	if len(t.DependsOnTaskIDs) > 0 { //&& t.ParentID == nil {
		var existingTasks []entities.Task
		if err := r.db.Where("id IN ?", t.DependsOnTaskIDs).Find(&existingTasks).Error; err != nil {
			return task.ErrFailedToFindDependsOnTasks
		}

		if len(existingTasks) != len(t.DependsOnTaskIDs) {
			return task.ErrFailedToFindDependsOnTasks
		}

		var taskDependencies []entities.TaskDependency
		for _, dependencyID := range t.DependsOnTaskIDs {
			taskDependencies = append(taskDependencies, entities.TaskDependency{
				DependentTaskID:  taskEntity.ID,
				DependencyTaskID: dependencyID,
			})
		}
		if err := r.db.Create(&taskDependencies).Error; err != nil {
			return task.ErrFailedToCreateTaskDependencies
		}

	}
	return nil
}

func (r *taskRepo) AddDependency(ctx context.Context, t *task.Task) error {
	//check if depends on task ID exist
	if len(t.DependsOnTaskIDs) > 0 {
		var errs error

		var existingTasks []entities.Task
		if err := r.db.Where("id IN ?", t.DependsOnTaskIDs).Find(&existingTasks).Error; err != nil {
			return task.ErrFailedToFindDependsOnTasks
		}

		if len(existingTasks) != len(t.DependsOnTaskIDs) {
			return task.ErrFailedToFindDependsOnTasks
		}
		// check dependency circle
		for _, tID := range t.DependsOnTaskIDs {
			if r.CheckCircularDependency(t.ID, tID) {
				errs = errors.Join(errs, fmt.Errorf("in circular dependency with task %v", tID))
			} else {
				if err := r.db.WithContext(ctx).Model(&entities.Task{ID: t.ID}).
					Association("DependsOn").
					Append(&entities.Task{ID: tID}); err != nil {
					return errors.Join(errs, err)
				}
			}
		}
	}
	return nil
}

// using dfs

func (r *taskRepo) CheckCircularDependency(taskID, dependencyID uuid.UUID) bool {
	visited := make(map[uuid.UUID]bool)
	var dfs func(current uuid.UUID) bool

	dfs = func(current uuid.UUID) bool {
		if current == taskID {
			return true
		}
		if visited[current] {
			return false
		}
		visited[current] = true

		var dependencies []uuid.UUID
		r.db.Model(&entities.Task{}).Where("id = ?", current).Association("DependsOn").Find(&dependencies)

		for _, depID := range dependencies {
			if dfs(depID) {
				return true
			}
		}
		return false
	}
	dependentCircle := dfs(dependencyID)
	return dependentCircle
}

func (r *taskRepo) GetBoardID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	var t entities.Task

	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", id).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t.BoardID, nil
}
