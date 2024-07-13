package presenter

import (
	"server/internal/notification"
	"time"

	"server/pkg/fp"

	"github.com/google/uuid"
)

type NotifResp struct {
	CreatedAt        time.Time                     `json:"created_at"`
	ID               uuid.UUID                     `json:"id"`
	IsSeen           bool                          `json:"is_seen"`
	Description      string                        `json:"desc"`
	NotificationType notification.NotificationType `json:"notif_type"`
}

func DomainNotifToNotifResp(n notification.Notification) NotifResp {
	return NotifResp{
		CreatedAt:        n.CreatedAt,
		ID:               n.ID,
		IsSeen:           n.IsSeen,
		Description:      n.Description,
		NotificationType: n.NotificationType,
	}
}

func BatchNotifToNotifResp(n []notification.Notification) []NotifResp {

	return fp.Map(n, DomainNotifToNotifResp)
}
