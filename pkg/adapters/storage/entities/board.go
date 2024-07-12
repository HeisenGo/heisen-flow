package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `gorm:"index"`
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Relationships
	//Users     []User `gorm:"many2many:user_board_roles;constraint:OnDelete:CASCADE;"`
	Tasks      []Task `gorm:"foreignKey:BoardID"`
	UserBoardRoles []UserBoardRole `gorm:"foreignKey:BoardID"`
}
