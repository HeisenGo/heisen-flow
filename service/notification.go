package service

import (
	"context"
	"errors"
	"fmt"
	"server/internal/board"
	"server/internal/notification"
	u "server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/rbac"
	"github.com/google/uuid"
)

type NotificationService struct {
	userOps          *u.Ops
	boardOps         *board.Ops
	userBoardRoleOps *userboardrole.Ops
	notificationOps  *notification.Ops
}

func NewNotificationService(userOps *u.Ops, boardOps *board.Ops, userBoardOps *userboardrole.Ops) *NotificationService {
	return &NotificationService{userOps: userOps, boardOps: boardOps, userBoardRoleOps: userBoardOps}
}

func (s *NotificationService) GetNotifications(ctx context.Context, userID , boardID uuid.UUID) ([]notification.Notification, error) {
	user, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, u.ErrUserNotFound
	}

	return s.notificationOps.DisplyNotification(ctx, userID,boardID)
}