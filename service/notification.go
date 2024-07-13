package service

import (
	"context"
	"server/internal/notification"
	u "server/internal/user"

	"github.com/google/uuid"
)

type NotificationService struct {
	userOps         *u.Ops
	notificationOps *notification.Ops
}

func NewNotificationService(notificationOps *notification.Ops, userOps *u.Ops) *NotificationService {
	return &NotificationService{notificationOps: notificationOps,
	userOps: userOps,}
}

func (s *NotificationService) CreateNotification(ctx context.Context, n *notification.Notification) error {
	err := s.notificationOps.CreateNotification(ctx, n)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error) {
	user, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, u.ErrUserNotFound
	}

	return s.notificationOps.GetUserUnseenNotifications(ctx, userID)
}

func (s *NotificationService) MarkNotificationAsSeen(ctx context.Context, notificationID, userID uuid.UUID) (*notification.Notification, error) {
	user, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, u.ErrUserNotFound
	}
	notif, err := s.notificationOps.GetNotificationByID(ctx, notificationID)
	if err != nil {
		return nil, err
	}
	if notif.UserBoardRole.Role == "" {
		return nil, ErrPermissionDenied
	}
	notiff, err := s.notificationOps.MarkNotificationAsSeen(ctx, notificationID)
	if err != nil {
		return nil, err
	}
	return notiff, err
}
