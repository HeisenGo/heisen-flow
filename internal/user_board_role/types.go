package userboardrole

import (
	"context"
	"errors"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

var (
	ErrUserRoleNotFound = errors.New("user role not found")
	ErrWrongRole        = errors.New("wrong role")
)

type Repo interface {
	GetUserBoardRole(ctx context.Context, userID, boardID string) (rbac.Role, error)
	SetUserBoardRole(ctx context.Context, ub *UserBoardRole) error
	RemoveUserBoardRole(ctx context.Context, userID, boardID string) error
}

type UserBoardRole struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	BoardID uuid.UUID
	Role    rbac.Role
}
