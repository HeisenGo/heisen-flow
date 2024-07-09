package notification

import (
	"context"
	"github.com/google/uuid"
)

type Repo interface {
	CreateNotification(ctx context.Context, userID, boardID uuid.UUID) error
	DisplyNotification(ctx context.Context, userID, boardID uuid.UUID) (*Notification,error)
	DeleteNotification(ctx context.Context, notif *Notification) error
}

type Notification struct {
    ID                uuid.UUID
    ISSeen            bool
    Description       string
    NotificationType  string
    UserBoardRoleID   uuid.UUID `gorm:"type:uuid;not null"`
}