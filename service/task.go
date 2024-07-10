package service

import (
	"context"
	"fmt"
	b "server/internal/board"
	"server/internal/task"
	u "server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

// TaskService handles task-related operations
type TaskService struct {
	userOps          *u.Ops
	boardOps         *b.Ops
	userBoardRoleOps *userboardrole.Ops
	taskOps          *task.Ops
}

// NewTaskService creates a new BoardService
func NewTaskService(userOps *u.Ops, boardOps *b.Ops, userBoardOps *userboardrole.Ops, taskOps *task.Ops) *TaskService {
	return &TaskService{userOps: userOps,
		boardOps:         boardOps,
		userBoardRoleOps: userBoardOps,
		taskOps:          taskOps}
}

func (s *BoardService) GetUserTasks(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]task.Task, uint, error) {
	return nil, 0, nil
}

func (s *TaskService) CreateTask(ctx context.Context, task *task.Task) error {
	// check if the creator exists
	user, err := s.userOps.GetUserByID(ctx, task.CreatedByUserID)
	if err != nil {
		return err
	}

	if user == nil {
		return u.ErrUserNotFound
	}

	// check if the board exists
	board, err := s.boardOps.GetBoardByID(ctx, task.BoardID)
	if err != nil {
		return err
	}

	if board == nil {
		return b.ErrBoardNotFound
	}

	// To Do
	//check if parent id is not null the parent task exists
	if task.ParentID!=nil{
		fmt.Print("Not implemented")
	}

	// check if assignee exists in this board
	if task.AssigneeUserID != nil {
		// assignee, err := s.userOps.GetUserByID(ctx, *task.AssigneeUserID)
		// if err != nil {
		// 	return err
		// }

		// if assignee == nil {
		// 	return errors.New("assignee not found")
		// }
		// get role ? can viewer be assigned a task???
		// check membership if assignee is not empty
		role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, *task.AssigneeUserID, board.ID)
		if err != nil {
			return err
		}

		if role == "" {
			return ErrNotMember
		}
	}

	// check permission for creator
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, user.ID, board.ID)
	if err != nil {
		return err
	}

	if !rbac.HasPermission(role, rbac.PermissionCreateTask) {
		return ErrPermissionDenied
	}

	err = s.taskOps.Create(ctx, task)
	if err != nil {
		return err
	}
	return nil
}
