package mappers

import (
	userboardrole "server/internal/user_board_role"
	"server/pkg/adapters/storage/entities"
)

func UserBoardRoleDomainToEntity(b *userboardrole.UserBoardRole) *entities.UserBoardRole {
	return &entities.UserBoardRole{
		UserID:   b.UserID,
		BoardID:  b.BoardID,
		UserRole: string(b.Role),
	}
}
