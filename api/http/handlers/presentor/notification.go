package presenter

import (
	"server/internal/notification"

	"github.com/google/uuid"
)

type NotificationReq struct {
	UserBoardRoleID               uuid.UUID `json:"user_board_role_id"`
	Description      string    `json:"description"`
	ISSeen           bool      `json:"is_seen"`
	NotificationType string    `json:"type"`
}

func NotificationToNotificationDomain(not *NotificationReq) *notification.Notification {
	return &notification.Notification{
		UserBoardRoleID: not.UserBoardRoleID,
		Description:not.Description  ,
		NotificationType:     not.NotificationType,
		ISSeen:  not.ISSeen,
	}
}
