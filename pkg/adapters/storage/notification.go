package storage

import (
	"context"
	"server/internal/notification"
	"server/pkg/adapters/storage/entities"
	"server/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) notification.Repo {
	return &notificationRepo{
		db: db,
	}
}

func (r *notificationRepo) CreateNotification(ctx context.Context, notif *notification.Notification) error {
    var userBoardRole entities.UserBoardRole
    if err := r.db.WithContext(ctx).First(&userBoardRole, "id = ?", notif.UserBoardRoleID).Error; err != nil {
        return err
    }

    newNotification := mappers.NotificationDomainToEntity(notif)
    if err := r.db.WithContext(ctx).Save(&newNotification).Error; err != nil {
        return err
    }

    notif.ID = newNotification.ID
    return nil
}

func (r *notificationRepo) GetUserUnseenNotifications(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error) {
    var notifications []entities.Notification

    if err := r.db.WithContext(ctx).
        Joins("JOIN user_board_roles ON user_board_roles.id = notifications.user_board_role_id").
        Where("user_board_roles.user_id = ? AND notifications.is_seen = ?", userID, false).
        Find(&notifications).Error; err != nil {
        return nil, err
    }
	var domainNotifications []notification.Notification
	for _, notif := range notifications {
        domainNotification := mappers.NotificationEntityToDomain(&notif)
        domainNotifications = append(domainNotifications, *domainNotification)
    }
    return domainNotifications, nil
}

func (r *notificationRepo) MarkNotificationAsSeen(ctx context.Context, notificationID uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.Notification{}).
        Where("id = ?", notificationID).
        Update("is_seen", true).Error
}