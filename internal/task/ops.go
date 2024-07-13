package task

import (
	"context"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) GetTaskByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	task, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (o *Ops) GetFullTaskByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	task, err := o.repo.GetFullByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}
func (o *Ops) UpdateTaskColumnByID(ctx context.Context, taskID uuid.UUID, colID uuid.UUID) (*Task, error) {
	updatedTask, err := o.repo.UpdateTaskColumnByID(ctx, taskID, colID)
	if err != nil {
		return nil, err
	}
	return updatedTask, nil
}

func (o *Ops) Create(ctx context.Context, task *Task) error {

	if err := validateTitleAndDescription(task.Title, task.Description); err != nil {
		return err
	}
	if task.StoryPoint != uint(0) {
		if err := validateStoryPoint(task.StoryPoint); err != nil {
			return err
		}
	}
	return o.repo.Insert(ctx, task)
}

func (o *Ops) AddDependency(ctx context.Context, t *Task) error {
	return o.repo.AddDependency(ctx, t)
}

func (o *Ops) ReorderTasks(ctx context.Context, colID uuid.UUID, newOrder map[uuid.UUID]uint) ([]Task, error) {
	return o.repo.ReorderTasks(ctx, colID, newOrder)
}
