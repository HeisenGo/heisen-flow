package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relationships
	TaskID uuid.UUID `gorm:"type:uuid;not null"`
	Task   *Task     `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserBoardRoleID uuid.UUID      `gorm:"type:uuid;not null"`
	UserBoardRole   *UserBoardRole `gorm:"foreignKey:UserBoardRoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
