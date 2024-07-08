package entities

import (
	"server/pkg/rbac"

	"gorm.io/gorm"
)

// UserBoardRole struct: Represents the relationship between a user, a board, and the user's role on that board.
type UserBoardRole struct {
	gorm.Model // to do use UUID
	UserID     string
	BoardID    string
	Role       rbac.Role
}
