package notification

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Repo interface {
	CreateNotification(ctx context.Context, notif *Notification) error
	GetUserUnseenNotifications(ctx context.Context, userID uuid.UUID) ([]Notification, error)
	MarkNotificationAsSeen(ctx context.Context, notificationID uuid.UUID) error
}

type Notification struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID                uuid.UUID
    IsSeen            bool
    Description       string
    NotificationType  string
    UserBoardRoleID   uuid.UUID `gorm:"type:uuid;not null"`
}