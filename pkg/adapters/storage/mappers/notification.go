package mappers

import (
	"server/internal/notification"
	"server/pkg/adapters/storage/entities"
)

func NotificationEntityToDomain(entity *entities.Notification) *notification.Notification {
	return &notification.Notification{
		ID:               entity.ID,
		IsSeen:           entity.IsSeen,
		Description:      entity.Description,
		NotificationType: entity.NotificationType,
		UserBoardRoleID: entity.UserBoardRoleID,
	}
}

func NotificationDomainToEntity(domainNotification *notification.Notification) *entities.Notification {
	return &entities.Notification{
		IsSeen:           domainNotification.IsSeen,
		Description:      domainNotification.Description,
		NotificationType: domainNotification.NotificationType,
		UserBoardRoleID:  domainNotification.UserBoardRoleID,
	}
}
