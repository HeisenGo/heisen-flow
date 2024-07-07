package userboard

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

func (o *Ops) SetUserBoardRole(ctx context.Context, userID, boardID string, role rbac.Role) error {
	return o.repo.SetUserBoardRole(ctx, userID, boardID, role)
}

func (o *Ops) RemoveUserBoardRole(ctx context.Context, userID, boardID string) error {
	return o.repo.RemoveUserBoardRole(ctx, userID, boardID)
}
