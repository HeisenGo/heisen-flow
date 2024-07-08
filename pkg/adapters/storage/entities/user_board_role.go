package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserBoardRole struct: Represents the relationship between a user, a board, and the user's role on that board.
type UserBoardRole struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID      `gorm:"index"`
	User      *User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BoardID   uuid.UUID      `gorm:"index"`
	Board     *Board         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserRole  string
}
