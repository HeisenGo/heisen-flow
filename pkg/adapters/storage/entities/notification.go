package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	IsSeen           bool
	Description      string
	NotificationType string
	UserBoardRoleID  uuid.UUID `gorm:"type:uuid"` //Assignee
}
