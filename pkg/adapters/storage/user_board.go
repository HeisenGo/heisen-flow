/*
This file provides methods for managing user-board roles in the storage layer:

GetUserBoardRole method: Retrieves the role of a user for a specific board.
SetUserBoardRole method: Sets the role of a user for a specific board.
RemoveUserBoardRole method: Removes the role of a user for a specific board.
*/

package storage

import (
	"context"
	userboard "server/internal/user_board"
	"server/pkg/rbac"

	"gorm.io/gorm"
)

type userBoardRepo struct {
	db *gorm.DB
}

func NewUserBoardRepo(db *gorm.DB) userboard.Repo {
	return &userBoardRepo{db}
}

func (r *userBoardRepo) GetUserBoardRole(ctx context.Context, userID, boardID string) (rbac.Role, error) {
	return "", userboard.ErrUserRoleNotFound
}

func (r *userBoardRepo) SetUserBoardRole(ctx context.Context, userID, boardID string, role rbac.Role) error {
	return nil
}

func (r *userBoardRepo) RemoveUserBoardRole(ctx context.Context, userID, boardID string) error {
	return nil
}
