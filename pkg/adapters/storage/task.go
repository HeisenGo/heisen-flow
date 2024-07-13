package storage

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
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

	var maxTaskColumnOrder uint
	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("column_id = ?", t.ColumnID).Select("COALESCE(MAX(\"order\"), 0)").Scan(&maxTaskColumnOrder).Error
	if err != nil {
		return err
	}
	taskEntity.Order = maxTaskColumnOrder
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

	return nil
}

// CheckCircularDependency using Depth First Search
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

func (r *taskRepo) GetFullByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	var t entities.Task

	if err := r.db.
		Preload("UserBoardRole").
		Preload("UserBoardRole.User").
		Preload("Parent").
		Preload("Subtasks").
		Preload("Column").
		Preload("Board").
		Preload("DependsOn").
		//Preload("TODO:Comments").
		First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}
	domainTask := mappers.TaskEntityToDomain(t)
	return &domainTask, nil
}

func (r *taskRepo) UpdateTaskColumnByID(ctx context.Context, taskID uuid.UUID, colID uuid.UUID) (*task.Task, error) {
	var t entities.Task

	// Find the task by ID
	if err := r.db.WithContext(ctx).First(&t, "id = ?", taskID).Error; err != nil {
		return nil, task.ErrTaskNotFound
	}

	// Load the new column to check its name
	var newColumn entities.Column
	if err := r.db.WithContext(ctx).First(&newColumn, "id = ?", colID).Error; err != nil {
		return nil, task.ErrColumnNotFound
	}
	// If the new column name is "done", delete relevant TaskDependency records
	if newColumn.Name == "done" {
		// Check if there are any dependencies where this task is a dependent
		var dependencyCount int64
		if err := r.db.WithContext(ctx).Model(&entities.TaskDependency{}).Where("dependent_task_id = ?", t.ID).Count(&dependencyCount).Error; err != nil {
			return nil, err
		}

		// If there are dependencies, abort the update
		if dependencyCount > 0 {
			return nil, task.ErrCantDoneDependentTask
		}
		// Update the task's ColumnID
		t.ColumnID = colID
		if err := r.db.WithContext(ctx).Where("dependency_task_id = ?", t.ID).Delete(&entities.TaskDependency{}).Error; err != nil {
			return nil, err
		}
	}

	// Save all subtasks using Association
	if err := r.db.Model(&entities.Task{}).Where("parent_id = ?", t.ID).Update("column_id", colID).Error; err != nil {
		return nil, err
	}

	t.ColumnID = colID
	// Save the updated task and its subtasks
	if err := r.db.WithContext(ctx).Save(&t).Error; err != nil {
		return nil, err
	}

	// Convert to domain entity if needed
	domainTask := mappers.TaskEntityToDomain(t)
	return &domainTask, nil
}

func (r *taskRepo) ReorderTasks(ctx context.Context, colID uuid.UUID, newOrder map[uuid.UUID]uint) ([]task.Task, error) {
	var tasks []entities.Task
	if err := r.db.WithContext(ctx).Where("column_id = ?", colID).Find(&tasks).Error; err != nil {
		return nil, task.ErrFailedToFetchTasks
	}
	if len(newOrder) != len(tasks) {
		return nil, task.ErrLengthMismatch
	}
	for taskID := range newOrder {
		found := false
		for _, t := range tasks {
			if t.ID == taskID {
				found = true
				break
			}
		}
		if !found {
			return nil, task.ErrInvalidTaskID
		}
	}

	var maxOrder uint
	for _, t := range tasks {
		if t.Order > maxOrder {
			maxOrder = t.Order
		}
	}
	tempOrder := maxOrder + 1

	for _, t := range tasks {
		if err := r.db.WithContext(ctx).Model(&t).Update("order", tempOrder).Error; err != nil {
			return nil, task.ErrFailedToUpdateTask
		}
		tempOrder++
	}

	for _, t := range tasks {
		newOrderNum, exists := newOrder[t.ID]
		if !exists {
			continue
		}
		if err := r.db.WithContext(ctx).Model(&t).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "title"}, {Name: "order"}}}).Update("order", newOrderNum).Error; err != nil {
			return nil, task.ErrFailedToUpdateTask
		}
		t.Order = newOrderNum
	}

	domainTasks := mappers.BatchTaskEntitiesToDomain(tasks)
	return domainTasks, nil

}
