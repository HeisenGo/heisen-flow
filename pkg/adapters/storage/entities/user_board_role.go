package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserBoardRole struct: Represents the relationship between a user, a board, and the user's role on that board.
type UserBoardRole struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	BoardID   uuid.UUID `gorm:"type:uuid;not null"`
	UserRole  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User  *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Board *Board `gorm:"foreignKey:BoardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
