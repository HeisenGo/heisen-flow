package service

import (
	"context"
	"server/internal/notification"
	u "server/internal/user"
	"github.com/google/uuid"
)

type NotificationService struct {
	userOps          *u.Ops
	notificationOps  *notification.Ops
}

func NewNotificationService(notificationOps *notification.Ops) *NotificationService {
	return &NotificationService{notificationOps: notificationOps}
}

func (s *NotificationService) CreateNotification(ctx context.Context,n *notification.Notification) error{
	err := s.notificationOps.CreateNotification(ctx,n)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) GetUserUnseenNotifications(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error) {
	user, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, u.ErrUserNotFound
	}

	return s.notificationOps.GetUserUnseenNotifications(ctx, userID)
}

func (s *NotificationService) MarkNotificationAsSeen(ctx context.Context, notificationID uuid.UUID)  error {
	err := s.notificationOps.MarkNotificationAsSeen(ctx,notificationID)
	return err
}