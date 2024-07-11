package storage

import (
	"context"
	"server/internal/notification"
	"server/pkg/adapters/storage/mappers"
	"server/pkg/adapters/storage/entities"
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

// func (r *notificationRepo) DisplyNotification(ctx context.Context, userID, boardID uuid.UUID) ([]notification.Notification,error) {
// 	var n entities.Notification
// 	var notifs []notification.Notification
// 	err := r.db.WithContext(ctx).Model(&entities.Notification{}).Where("id = ?", userID).First(&n).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
	
// 	append()
// 	return mappers.NotificationEntityToDomain(&n), nil
// }

// func (r *notificationRepo) DeleteNotification(ctx context.Context, notif *notification.Notification) error {

// }