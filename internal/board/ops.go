package board

import (
	"context"
	"time"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) UserBoards(ctx context.Context, userID uint, page, pageSize uint) ([]Board, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	return o.repo.GetUserBoards(ctx, userID, limit, offset)
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
