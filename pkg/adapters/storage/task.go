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
	if len(t.DependsOnTaskIDs) > 0 {
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
	if len(t.DependsOnTaskIDs) > 0 {
		var errs error

		var existingTasks []entities.Task
		if err := r.db.Where("id IN ?", t.DependsOnTaskIDs).Find(&existingTasks).Error; err != nil {
			return err
		}

		if len(existingTasks) != len(t.DependsOnTaskIDs) {
			return task.ErrFailedToFindDependsOnTasks
		}

		// Check for circular dependencies
		for _, dependsOnID := range t.DependsOnTaskIDs {
			if r.CheckCircularDependency(t.ID, dependsOnID) {
				errs = errors.Join(errs, fmt.Errorf("circular dependency detected with task %v", dependsOnID))
			}
		}
		if errs != nil {
			return errors.Join(task.ErrCircularDependency, errs)
		}
		// Retrieve the main task entity
		var tEntity entities.Task
		fmt.Println(t.ID)
		if err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", t.ID).First(&tEntity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return task.ErrTaskNotFound
			}
			return err
		}
		// TO DO : Bulk
		for _, existingTask := range existingTasks {
			if err := r.db.WithContext(ctx).Model(&tEntity).Association("DependsOn").Append(&existingTask); err != nil {
				return fmt.Errorf("failed to add dependency %v: %w", existingTask.ID, err)
			}
		}

	}
	return nil
}

// using Depth First Search
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

		var dependencies []entities.Task
		r.db.Model(&entities.Task{ID: current}).Association("DependsOn").Find(&dependencies)

		for _, dep := range dependencies {
			if dfs(dep.ID) {
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
