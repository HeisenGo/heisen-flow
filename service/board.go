package service

import (
	"context"
	"errors"
	"fmt"
	userboard "server/internal/user_board"
	"server/pkg/rbac"
)

// BoardService handles board-related operations
type BoardService struct {
	userBoardOps *userboard.Ops
}

// NewBoardService creates a new BoardService
func NewBoardService(userBoardOps *userboard.Ops) *BoardService {
	return &BoardService{userBoardOps: userBoardOps}
}

// CreateTask creates a new task in a board
func (s *BoardService) CreateTask(ctx context.Context, userID, boardID string, taskDetails map[string]interface{}) error {
	role, err := s.userBoardOps.GetUserBoardRole(ctx, userID, boardID)
	if err != nil {
		return err
	}

	if !rbac.HasPermission(role, rbac.PermissionCreateTask) {
		return errors.New("permission denied: cannot create task")
	}

	// To Do create task
	return nil
}

// MoveTask moves a task to a different column
func (s *BoardService) MoveTask(ctx context.Context, userID, boardID, taskID, newColumnID string) error {
	role, err := s.userBoardOps.GetUserBoardRole(ctx, userID, boardID)
	if err != nil {
		return err
	}
	fmt.Print(role)
	// to do
	// task, err := s.taskOps.GetTask(taskID)
	// if err != nil {
	// 	return err
	// }

	// if task.AssigneeID == userID {
	// 	if !rbac.HasPermission(role, rbac.PermissionMoveOwnTask) {
	// 		return errors.New("permission denied: cannot move own task")
	// 	}
	// } else {
	// 	if !rbac.HasPermission(role, rbac.PermissionMoveAnyTask) {
	// 		return errors.New("permission denied: cannot move other's task")
	// 	}
	// }

	// Implement task moving

	return nil
}

// AddColumn adds a new column to the board
func (s *BoardService) AddColumn(ctx context.Context, userID, boardID string, columnDetails map[string]interface{}) error {
	role, err := s.userBoardOps.GetUserBoardRole(ctx, userID, boardID)
	if err != nil {
		return err
	}

	if !rbac.HasPermission(role, rbac.PermissionManageColumns) {
		return errors.New("permission denied: cannot add column")
	}

	// column addition logic

	return nil
}

// InviteUser invites a user to the board
func (s *BoardService) InviteUser(ctx context.Context, inviterID, boardID, inviteeID string, role rbac.Role) error {
	inviterRole, err := s.userBoardOps.GetUserBoardRole(ctx, inviterID, boardID)
	if err != nil {
		return err
	}

	if !rbac.HasPermission(inviterRole, rbac.PermissionInviteUsers) {
		return errors.New("permission denied: cannot invite users")
	}

	// Implement user invitation

	return nil
}
