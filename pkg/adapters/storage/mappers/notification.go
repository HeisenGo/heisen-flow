package mappers

import (
	"server/internal/notification"
	"server/pkg/adapters/storage/entities"
)

func NotificationEntityToDomain(entity *entities.Notification) *notification.Notification {
	return &notification.Notification{
		ID:               entity.ID,
		ISSeen:           entity.ISSeen,
		Description:      entity.Description,
		NotificationType: entity.NotificationType,
	}
}

func NotificationDomainToEntity(domainNotification *notification.Notification) *entities.Notification {
	return &entities.Notification{
		ID:               domainNotification.ID,
		ISSeen:           domainNotification.ISSeen,
		Description:      domainNotification.Description,
		NotificationType: domainNotification.NotificationType,
	}
}
