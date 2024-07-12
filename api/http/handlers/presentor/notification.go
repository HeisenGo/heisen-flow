package presenter

import (
	"server/internal/notification"
	"time"

	"github.com/google/uuid"
)

type NotificationReq struct {
	UserBoardRoleID               uuid.UUID `json:"user_board_role_id"`
	Description      string    `json:"description"`
	IsSeen           bool      `json:"is_seen"`
	NotificationType string    `json:"type"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NotificationToNotificationDomain(not *NotificationReq) *notification.Notification {
	return &notification.Notification{
		UserBoardRoleID: not.UserBoardRoleID,
		Description:not.Description  ,
		NotificationType:     not.NotificationType,
		IsSeen:  not.IsSeen,
	}
}
