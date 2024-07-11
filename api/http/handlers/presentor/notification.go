package presenter

import (
	"server/internal/notification"

	"github.com/google/uuid"
)

type NotificationReq struct {
	ID               uuid.UUID `json:"board_id"`
	Description      string    `json:"description"`
	ISSeen           bool      `json:"is_seen"`
	NotificationType string    `json:"type"`
}

func NotificationToNotificationDomain(not *NotificationReq) *notification.Notification {
	return &notification.Notification{
		ID: not.ID,
		Description:not.Description  ,
		NotificationType:     not.NotificationType,
		ISSeen:  not.ISSeen,
	}
}
