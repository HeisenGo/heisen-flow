package userboardrole

import (
	"context"
	"server/pkg/rbac"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) GetUserBoardRole(ctx context.Context, userID, boardID string) (rbac.Role, error) {
	return o.repo.GetUserBoardRole(ctx, userID, boardID)
}

func (o *Ops) SetUserBoardRole(ctx context.Context, ub *UserBoardRole) error {
	return o.repo.SetUserBoardRole(ctx, ub)
}

func (o *Ops) RemoveUserBoardRole(ctx context.Context, userID, boardID string) error {
	return o.repo.RemoveUserBoardRole(ctx, userID, boardID)
}
