package userboard

import (
	"context"
	"errors"
	"server/pkg/rbac"
)

var (
	ErrUserRoleNotFound = errors.New("user role not found")
	ErrWrongRole        = errors.New("wrong role")
)

type Repo interface {
	GetUserBoardRole(ctx context.Context, userID, boardID string) (rbac.Role, error)
	SetUserBoardRole(ctx context.Context, userID, boardID string, role rbac.Role) error
	RemoveUserBoardRole(ctx context.Context, userID, boardID string) error
}

type UserBoard struct {
	ID      uint
	UserID  string
	BoardID string
	Role    rbac.Role
}
