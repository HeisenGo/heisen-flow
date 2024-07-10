package storage

import (
	"context"
	"errors"
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

func (r *notificationRepo) CreateNotification(ctx context.Context, userID, boardID uuid.UUID) error {
	var n notification.Notification
	newNotification := mappers.NotificationDomainToEntity(&n)
	err := r.db.WithContext(ctx).Create(&newNotification).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *notificationRepo) DisplyNotification(ctx context.Context, userID, boardID uuid.UUID) ([]notification.Notification,error) {

}

func (r *notificationRepo) DeleteNotification(ctx context.Context, notif *notification.Notification) error {

}