package notification

import (
	"context"
	"errors"
	"server/internal/task"
	userboardrole "server/internal/user_board_role"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	UserInvited     = NotificationType("Invite User")
	TaskMoved       = NotificationType("Move Task")
	CommentedNotif  = NotificationType("Comment")
	TaskUpdateNotif = NotificationType("Update Task")
)

var (
	ErrFailedToCreateNotif = errors.New("Failed to create notif")
)

var (
	ErrNotifNotFound  = errors.New("notif not found")
	ErrNotifsNotFound = errors.New("notifications not found")
)

type Repo interface {
	CreateNotification(ctx context.Context, notif *Notification) error
	GetUserUnseenNotifications(ctx context.Context, userID uuid.UUID) ([]Notification, error)
	MarkNotificationAsSeen(ctx context.Context, notificationID uuid.UUID) (*Notification, error)
	GetNotificationByID(ctx context.Context, notificationID uuid.UUID) (*Notification, error)
	NotifBroadCasting(ctx context.Context, notif *Notification, boardID, userID uuid.UUID, task *task.Task)error
}

type Notification struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ID               uuid.UUID
	IsSeen           bool
	Description      string
	NotificationType NotificationType
	UserBoardRoleID  uuid.UUID `gorm:"type:uuid;not null"`
	UserBoardRole    *userboardrole.UserBoardRole
	BoardName        string
}

func NewNotification(description string, notificationType NotificationType, userBoardRoleID uuid.UUID) *Notification {
	return &Notification{
		IsSeen:           false,
		Description:      description,
		NotificationType: notificationType,
		UserBoardRoleID:  userBoardRoleID,
	}
}
