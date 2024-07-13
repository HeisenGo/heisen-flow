package mappers

import (
	"server/internal/notification"
	"server/pkg/adapters/storage/entities"
)

func NotificationEntityToDomain(entity *entities.Notification) *notification.Notification {

	return &notification.Notification{
		CreatedAt:        entity.CreatedAt,
		ID:               entity.ID,
		IsSeen:           entity.IsSeen,
		Description:      entity.Description,
		NotificationType: notification.NotificationType(entity.NotificationType),
		UserBoardRoleID:  entity.UserBoardRoleID,
	}
}

func NotificationDomainToEntity(domainNotification *notification.Notification) *entities.Notification {
	return &entities.Notification{
		IsSeen:           domainNotification.IsSeen,
		Description:      domainNotification.Description,
		NotificationType: string(domainNotification.NotificationType),
		UserBoardRoleID:  domainNotification.UserBoardRoleID,
		CreatedAt:        domainNotification.CreatedAt,
		UpdatedAt:        domainNotification.UpdatedAt,
	}
}
