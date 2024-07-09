package entities

import (
	"github.com/google/uuid"
)

type Notification struct {
    ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    ISSeen            bool
    Description       string
    NotificationType  string
    UserBoardRoleID   uuid.UUID `gorm:"type:uuid;not null"`
}
