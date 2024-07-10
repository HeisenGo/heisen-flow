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

// func (o *Ops) GetBoardByID(ctx context.Context, id uuid.UUID) (*Board, error) {
// 	board, err := o.repo.GetByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if board == nil {
// 		return nil, ErrBoardNotFound
// 	}
// 	return board, nil
// }

// func (o *Ops) UserBoards(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]Board, uint, error) {
// 	limit := pageSize
// 	offset := (page - 1) * pageSize

// 	return o.repo.GetUserBoards(ctx, userID, limit, offset)
// }

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

func (o *Ops) Create(ctx context.Context, task *Task) error {
	if err := validateTitleAndDescription(task.Title, task.Description); err != nil {
		return err
	}
	return o.repo.Insert(ctx, task)
}

func (o *Ops) AddDependency(ctx context.Context, t *Task) error {
	return o.repo.AddDependency(ctx, t)
}

func (o *Ops) GetBoardID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	boardID, err := o.repo.GetBoardID(ctx, id)
	if err != nil {
		return nil, err
	}
	if boardID == nil {
		return nil, ErrBoardNotFound
	}
	return boardID, nil
}
