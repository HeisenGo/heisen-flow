package service

import (
	"context"
	"errors"
	"fmt"
	"server/internal/board"
	"server/internal/column"
	u "server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

var (
	ErrPermissionDenied         = errors.New("permission denied")
	ErrNotMember                = errors.New("the assignee is not a member of this board")
	ErrCantAssigned             = errors.New("assignee is a viewer")
	ErrOwnerExists              = errors.New("owner already exists")
	ErrUndefinedRole            = errors.New("undefined role, role should be one of the following values:viewer, editor, maintainer")
	ErrPermissionDeniedToInvite = errors.New("permission denied: cannot invite users")
	ErrAMember                  = errors.New("user already is a member")
	ErrPermissionDeniedToDelete = errors.New("permission denied: can not delete the board")
)

// BoardService handles board-related operations
type BoardService struct {
	userOps          *u.Ops
	boardOps         *board.Ops
	userBoardRoleOps *userboardrole.Ops
	columnOps        *column.Ops
}

// NewBoardService creates a new BoardService
func NewBoardService(userOps *u.Ops, boardOps *board.Ops,
	userBoardOps *userboardrole.Ops,
	columnOps *column.Ops) *BoardService {
	return &BoardService{userOps: userOps,
		boardOps:         boardOps,
		userBoardRoleOps: userBoardOps,
		columnOps:        columnOps}
}

func (s *BoardService) GetFullBoardByID(ctx context.Context, userID uuid.UUID, boardID uuid.UUID) (*board.Board, error) {
	b, err := s.boardOps.GetFullBoardByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	if b.Type == "private" {
		fetcherRole, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, boardID)
		if err != nil {
			return nil, ErrPermissionDeniedToInvite
		}

		if !rbac.HasPermission(fetcherRole, rbac.PermissionViewBoard) {
			return nil, ErrPermissionDeniedToInvite
		}
	}
	return b, err
}

func (s *BoardService) GetUserBoards(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]board.Board, uint, error) {
	user, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, u.ErrUserNotFound
	}

	return s.boardOps.GetUserBoards(ctx, userID, page, pageSize)
}

func (s *BoardService) GetPublicBoards(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]board.Board, uint, error) {
	user, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, u.ErrUserNotFound
	}

	return s.boardOps.GetPublicBoards(ctx, userID, page, pageSize)
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
	ub.Role = string(rbac.RoleOwner)
	err = s.userBoardRoleOps.SetUserBoardRole(ctx, ub)
	if err != nil {
		return err
	}
	// set first "done" default column

	col, err := s.columnOps.SetDoneAsDefault(ctx, ub.BoardID)
	if err != nil {
		return err
	}
	b.Columns = append(b.Columns, *col)
	return nil
}

// MoveTask moves a task to a different column
func (s *BoardService) MoveTask(ctx context.Context, userID, boardID uuid.UUID, taskID, newColumnID string) error {
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

// InviteUser invites a user to the board
func (s *BoardService) InviteUser(ctx context.Context, inviterID uuid.UUID, inviteeEmail string, userBoardRole *userboardrole.UserBoardRole) error {

	if userBoardRole.Role == string(rbac.RoleOwner) {
		return ErrOwnerExists
	}
	isPossibleRole := rbac.IsAPossibleRole(userBoardRole.Role)

	if !isPossibleRole {
		return ErrUndefinedRole
	}

	inviterRole, err := s.userBoardRoleOps.GetUserBoardRole(ctx, inviterID, userBoardRole.BoardID)
	if err != nil {
		return ErrPermissionDeniedToInvite
	}

	if !rbac.HasPermission(inviterRole, rbac.PermissionInviteUsers) {
		return ErrPermissionDeniedToInvite
	}

	invitedUser, err := s.userOps.GetUserByEmail(ctx, inviteeEmail)
	if err != nil {
		return err
	}
	if invitedUser == nil {
		return u.ErrUserNotFound
	}
	userBoardRole.UserID = invitedUser.ID
	b, err := s.boardOps.GetBoardByID(ctx, userBoardRole.BoardID)
	if err != nil {
		return err
	}

	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, invitedUser.ID, b.ID)
	if role != "" && err == nil {
		return ErrAMember
	}

	err = s.userBoardRoleOps.SetUserBoardRole(ctx, userBoardRole)
	if err != nil {
		return err
	}
	// apply your notification record create here
	return nil
}

func (s *BoardService) DeleteBoardByID(ctx context.Context, ub *userboardrole.UserBoardRole) error {
	// check board exists
	b, err := s.boardOps.GetBoardByID(ctx, ub.BoardID)

	if err != nil {
		return err
	}
	if b == nil {
		return board.ErrBoardNotFound
	}

	// check user exists
	user, err := s.userOps.GetUserByID(ctx, ub.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return u.ErrUserNotFound
	}

	// check user has permission of removing
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, ub.UserID, ub.BoardID)
	if err != nil {
		return ErrPermissionDeniedToInvite
	}

	if !rbac.HasPermission(role, rbac.PermissionRemoveBoard) {
		return ErrPermissionDeniedToDelete
	}

	err = s.boardOps.Delete(ctx, ub.BoardID)
	if err != nil {
		return err
	}

	return nil
}
