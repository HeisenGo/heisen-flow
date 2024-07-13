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
	"server/pkg/adapters/storage/entities"
	"server/pkg/adapters/storage/mappers"
	"server/pkg/rbac"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userBoardRepo struct {
	db *gorm.DB
}

func NewUserBoardRepo(db *gorm.DB) userboardrole.Repo {
	return &userBoardRepo{db}
}

func (r *userBoardRepo) GetUserBoardRole(ctx context.Context, userID, boardID uuid.UUID) (rbac.Role, error) {
	var userBoardRole entities.UserBoardRole
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		First(&userBoardRole).Error
	if err != nil {
		return "", userboardrole.ErrUserRoleNotFound
	}
	return rbac.Role(userBoardRole.UserRole), nil
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

func (r *userBoardRepo) GetUserBoardRoleObj(ctx context.Context, userID, boardID uuid.UUID) (*userboardrole.UserBoardRole, error) {
	var userBoardRole entities.UserBoardRole
	err := r.db.
		Preload("User").
		Where("user_id = ? AND board_id = ?", userID, boardID).
		First(&userBoardRole).Error
	if err != nil {
		return nil, userboardrole.ErrUserRoleNotFound
	}
	ubrDomain := mappers.UserBoardRoleEntityToDomain(userBoardRole)
	return &ubrDomain, nil
}

func (r *userBoardRepo) GetUserIDByUserBoardRoleID(ctx context.Context, userBoardRoleID uuid.UUID) (*uuid.UUID, error) {
	var userBoardRole entities.UserBoardRole
	result := r.db.WithContext(ctx).First(&userBoardRole, "id = ?", userBoardRoleID)
	if result.Error != nil {
		return nil, userboardrole.ErrUserRoleNotFound
	}
	return &userBoardRole.UserID, nil
}
