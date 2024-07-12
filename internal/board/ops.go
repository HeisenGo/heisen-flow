package board

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) GetBoardByID(ctx context.Context, id uuid.UUID) (*Board, error) {
	board, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}
	return board, nil
}

func (o *Ops) GetUserBoards(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]Board, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	return o.repo.GetUserBoards(ctx, userID, limit, offset)
}

func (o *Ops) GetPublicBoards(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]Board, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	return o.repo.GetPublicBoards(ctx, userID, limit, offset)
}

func (o *Ops) Create(ctx context.Context, board *Board) error {
	if err := ValidateBoardName(board.Name); err != nil {
		return ErrInvalidName
	}

	if board.Type != string(Private) && board.Type != string(Public) {
		return ErrWrongType
	}
	if board.CreatedAt.After(time.Now()) {
		return ErrWrongBoardTime
	}

	return o.repo.Insert(ctx, board)
}
