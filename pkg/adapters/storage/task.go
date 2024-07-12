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
	var errs error
	// Retrieve the main task entity
	var tEntity entities.Task
	fmt.Println(t.ID)
	if err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", t.ID).First(&tEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return task.ErrTaskNotFound
		}
		return err
	}

	// for _, tdi := range t.DependsOnTaskIDs {
	// 	var count int64
	// 	//r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", t.ID).First(&tEntity).Error
	// 	err := r.db.WithContext(ctx).Model(&entities.TaskDependency{}).Where("dependent_task_id = ? AND dependency_task_id = ?", t.ID, tdi).Count(&count).Error
	// 	if err != nil {
	// 		return fmt.Errorf("failed to check existing dependency %v: %w", tdi, err)
	// 	}
	// 	if count > 0 {
	// 		return fmt.Errorf("duplicate dependency detected: task %v already depends on task %v", t.ID, tdi)
	// 	}
	// }

	existingTasks, err := r.GetExistingTasks(ctx, t.DependsOnTaskIDs)

	if len(existingTasks) != len(t.DependsOnTaskIDs) || err != nil {
		return task.ErrFailedToFindDependsOnTasks
	}

	for i, _ := range t.DependsOnTaskIDs {
		fmt.Println(existingTasks[i].BoardID, t.BoardID)
		if existingTasks[i].Board.ID != t.BoardID {
			return errors.Join(task.ErrConflictedBoards, fmt.Errorf(" with %v", existingTasks[i].ID))
		}
	}

	// Check for circular dependencies and being in the same board
	for _, dependsOnID := range t.DependsOnTaskIDs {
		if r.CheckCircularDependency(t.ID, dependsOnID) {
			errs = errors.Join(errs, fmt.Errorf("circular dependency detected with task %v", dependsOnID))
		}
	}
	if errs != nil {
		return errors.Join(task.ErrCircularDependency, errs)
	}

	// TO DO : Bulk
	//var tDE entities.TaskDependency
	for _, existingTask := range existingTasks {
		if err := r.db.WithContext(ctx).Model(&tEntity).Association("DependsOn").Append(&existingTask); err != nil {
			return fmt.Errorf("failed to add dependency %v: %w", existingTask.ID, err)
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

func (r *taskRepo) GetExistingTasks(ctx context.Context, dependsOnTaskIDs []uuid.UUID) ([]entities.Task, error) {
	var existingTasks []entities.Task
	if err := r.db.
		Where("id IN ?", dependsOnTaskIDs).
		Preload("Board").
		Find(&existingTasks).Error; err != nil {
		return nil, err
	}
	return existingTasks, nil
}
