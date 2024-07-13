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

func (r *notificationRepo) NotifBroadCasting(ctx context.Context, notif *notification.Notification, boardID, userID uuid.UUID) error {
	var userBoardRoles []entities.UserBoardRole
	err := r.db.Where("board_id = ? AND user_role IN ?", boardID, []string{"maintainer", "owner"}).
		Find(&userBoardRoles).Error
	if err != nil {
		return notification.ErrFailedToCreateNotif
	}

	for i,_ := range userBoardRoles {
		newNotification := mappers.NotificationDomainToEntity(notif)
		notif.UserBoardRoleID = userBoardRoles[i].ID
		if err := r.db.WithContext(ctx).Save(&newNotification).Error; err != nil {
			return notification.ErrFailedToCreateNotif
		}
	}
	return nil
}

func (r *notificationRepo) CreateNotification(ctx context.Context, notif *notification.Notification) error {
	var userBoardRole entities.UserBoardRole
	//	TODO REMOVE REPEATED ONES
	if err := r.db.WithContext(ctx).First(&userBoardRole, "id = ?", notif.UserBoardRoleID).Error; err != nil {
		return notification.ErrFailedToCreateNotif
	}

	newNotification := mappers.NotificationDomainToEntity(notif)
	if err := r.db.WithContext(ctx).Save(&newNotification).Error; err != nil {
		return notification.ErrFailedToCreateNotif
	}

	notif.ID = newNotification.ID
	return nil
}

func (r *notificationRepo) GetUserUnseenNotifications(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error) {
	var notifications []entities.Notification

	result := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Joins("LEFT JOIN user_board_roles ON notifications.user_board_role_id = user_board_roles.id").
		Where("user_board_roles.user_id = ?", userID).
		Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}
	var domainNotifications []notification.Notification
	for _, notif := range notifications {
		domainNotification := mappers.NotificationEntityToDomain(&notif)
		domainNotifications = append(domainNotifications, *domainNotification)
	}

	return domainNotifications, nil
}

func (r *notificationRepo) MarkNotificationAsSeen(ctx context.Context, notificationID uuid.UUID) (*notification.Notification, error) {
	notification := &entities.Notification{}

	result := r.db.WithContext(ctx).Model(&entities.Notification{}).
		Where("id = ?", notificationID).
		Update("is_seen", true).
		First(notification) // Retrieve the updated entity

	if result.Error != nil {
		return nil, result.Error
	}
	notif := mappers.NotificationEntityToDomain(notification)
	return notif, nil
}

func (r *notificationRepo) GetNotificationByID(ctx context.Context, notificationID uuid.UUID) (*notification.Notification, error) {
	notification := &entities.Notification{}

	result := r.db.WithContext(ctx).Preload("UserBoardRole").
		Where("id = ?", notificationID).
		First(notification)

	if result.Error != nil {
		return nil, result.Error
	}
	notif := mappers.NotificationEntityToDomain(notification)
	return notif, nil
}
