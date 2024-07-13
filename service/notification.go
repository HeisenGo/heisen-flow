package service

import (
	"context"
	"server/internal/notification"
	u "server/internal/user"
	userboardrole "server/internal/user_board_role"

	"github.com/google/uuid"
)

type NotificationService struct {
	userOps          *u.Ops
	notificationOps  *notification.Ops
	userBoardRoleOps *userboardrole.Ops
}

func NewNotificationService(notificationOps *notification.Ops, userOps *u.Ops, userBoardRoleOps *userboardrole.Ops) *NotificationService {
	return &NotificationService{notificationOps: notificationOps,
		userOps:          userOps,
		userBoardRoleOps: userBoardRoleOps}
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
	uID, err := s.userBoardRoleOps.GetUserIDByUserBoardRoleID(ctx, notif.UserBoardRoleID)
	if err != nil {
		return nil, err
	}
	if *uID != userID {
		return nil, ErrPermissionDenied
	}
	notiff, err := s.notificationOps.MarkNotificationAsSeen(ctx, notificationID)
	if err != nil {
		return nil, err
	}
	return notiff, err
}
