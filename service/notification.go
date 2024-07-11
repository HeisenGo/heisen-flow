package service

import (
	"context"
	"server/internal/board"
	"server/internal/notification"
	u "server/internal/user"
	userboardrole "server/internal/user_board_role"
)

type NotificationService struct {
	userOps          *u.Ops
	boardOps         *board.Ops
	userBoardRoleOps *userboardrole.Ops
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

// func (s *NotificationService) GetNotifications(ctx context.Context, userID , boardID uuid.UUID) ([]notification.Notification, error) {
// 	user, err := s.userOps.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if user == nil {
// 		return nil, u.ErrUserNotFound
// 	}

// 	return s.notificationOps.DisplyNotification(ctx, userID,boardID)
// }