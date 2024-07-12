package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title           string         `gorm:""`
	Description     string         `gorm:""`
	UserBoardRoleID uuid.UUID      `gorm:"type:uuid"`
	UserBoardRole   *UserBoardRole `gorm:"foreignKey:UserBoardRoleID"`
	TaskID          uuid.UUID      `gorm:"type:uuid"`
	Task            *Task          `gorm:"foreignKey:TaskID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
