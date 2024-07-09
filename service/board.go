package service

import (
	"context"
	"errors"
	"fmt"
	"server/internal/board"
	u "server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

// BoardService handles board-related operations
type BoardService struct {
	userOps          *u.Ops
	boardOps         *board.Ops
	userBoardRoleOps *userboardrole.Ops
}

// NewBoardService creates a new BoardService
func NewBoardService(userOps *u.Ops, boardOps *board.Ops, userBoardOps *userboardrole.Ops) *BoardService {
	return &BoardService{userOps: userOps, boardOps: boardOps, userBoardRoleOps: userBoardOps}
}

func (s *BoardService) GetUserBoards(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]board.Board, uint, error) {
	user, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, u.ErrUserNotFound
	}

	return s.boardOps.UserBoards(ctx, userID, page, pageSize)
}

func (s *BoardService) CreateBoard(ctx context.Context, b *board.Board, ub *userboardrole.UserBoardRole) error {
	user, err := s.userOps.GetUserByID(ctx, ub.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return u.ErrUserNotFound
	}

	err = s.boardOps.Create(ctx, b)
	if err != nil {
		return err
	}

	ub.BoardID = b.ID
	ub.Role = rbac.RoleOwner
	err = s.userBoardRoleOps.SetUserBoardRole(ctx, ub)
	if err != nil {
		return err
	}
	return nil
	// return errors.New("test error") -> for testing transaction
}

// CreateTask creates a new task in a board
func (s *BoardService) CreateTask(ctx context.Context, userID, boardID string, taskDetails map[string]interface{}) error {
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, boardID)
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
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, boardID)
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
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, boardID)
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
	inviterRole, err := s.userBoardRoleOps.GetUserBoardRole(ctx, inviterID, boardID)
	if err != nil {
		return err
	}

	if !rbac.HasPermission(inviterRole, rbac.PermissionInviteUsers) {
		return errors.New("permission denied: cannot invite users")
	}

	// Implement user invitation

	return nil
}
