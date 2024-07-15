package userboardrole

import (
	"context"
	"errors"
	"server/internal/user"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

var (
	ErrUserRoleNotFound = errors.New("user role not found")
	ErrWrongRole        = errors.New("wrong role")
)

type Repo interface {
	GetUserBoardRole(ctx context.Context, userID, boardID uuid.UUID) (rbac.Role, error)
	GetUserBoardRoleObj(ctx context.Context, userID, boardID uuid.UUID) (*UserBoardRole, error)
	SetUserBoardRole(ctx context.Context, ub *UserBoardRole) error
	RemoveUserBoardRole(ctx context.Context, userID, boardID string) error
	GetUserIDByUserBoardRoleID(ctx context.Context, userBoardRoleID uuid.UUID) (*uuid.UUID, error)
}

type UserBoardRole struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	User    *user.User
	BoardID uuid.UUID
	Role    string
}
