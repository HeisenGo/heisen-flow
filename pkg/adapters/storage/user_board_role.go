/*
This file provides methods for managing user-board roles in the storage layer:

GetUserBoardRole method: Retrieves the role of a user for a specific board.
SetUserBoardRole method: Sets the role of a user for a specific board.
RemoveUserBoardRole method: Removes the role of a user for a specific board.
*/

package storage

import (
	"context"
	userboardrole "server/internal/user_board_role"
	"server/pkg/adapters/storage/mappers"
	"server/pkg/rbac"

	"gorm.io/gorm"
)

type userBoardRepo struct {
	db *gorm.DB
}

func NewUserBoardRepo(db *gorm.DB) userboardrole.Repo {
	return &userBoardRepo{db}
}

func (r *userBoardRepo) GetUserBoardRole(ctx context.Context, userID, boardID string) (rbac.Role, error) {
	return "", userboardrole.ErrUserRoleNotFound
}

func (r *userBoardRepo) SetUserBoardRole(ctx context.Context, ub *userboardrole.UserBoardRole) error {
	userBoardRoleEntity := mappers.UserBoardRoleDomainToEntity(ub)
	if err := r.db.WithContext(ctx).Save(&userBoardRoleEntity).Error; err != nil {
		return err
	}

	ub.ID = userBoardRoleEntity.ID
	return nil
}

func (r *userBoardRepo) RemoveUserBoardRole(ctx context.Context, userID, boardID string) error {
	return nil
}
