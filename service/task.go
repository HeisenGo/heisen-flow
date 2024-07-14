package service

import (
	"context"
	"fmt"
	b "server/internal/board"
	"server/internal/column"
	"server/internal/notification"
	t "server/internal/task"
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
	taskOps          *t.Ops
	columnOps        *column.Ops
	notificaionOps   *notification.Ops
}

// NewTaskService creates a new TaskService
func NewTaskService(userOps *u.Ops, boardOps *b.Ops, userBoardOps *userboardrole.Ops, taskOps *t.Ops, columnOps *column.Ops, notifOps *notification.Ops) *TaskService {
	return &TaskService{userOps: userOps,
		boardOps:         boardOps,
		userBoardRoleOps: userBoardOps,
		taskOps:          taskOps,
		columnOps:        columnOps,
		notificaionOps:   notifOps,
	}
}

func (s *BoardService) GetUserTasks(ctx context.Context, userID uuid.UUID, page, pageSize uint) ([]t.Task, uint, error) {
	return nil, 0, nil
}

func (s *TaskService) CreateTask(ctx context.Context, task *t.Task) error {
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

	//check if parent id is not null and the parent task exists for sub tasks
	if task.ParentID != nil {
		_, err := s.taskOps.GetTaskByID(ctx, *task.ParentID)
		if err != nil {
			return t.ErrParentTaskNotFound
		}
	}

	// check if assignee exists in this board
	if task.AssigneeUserID != nil {

		// check membership if assignee is not empty
		role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, *task.AssigneeUserID, board.ID)
		if err != nil {
			return err
		}

		if role == "" {
			return ErrNotMember
		}
		// assignee can not be viewer
		if !rbac.HasPermission(role, rbac.PermissionMoveOwnTask) {
			return ErrCantAssigned
		}
		ubrObj, err := s.userBoardRoleOps.GetUserBoardRoleObj(ctx, *task.AssigneeUserID, board.ID)
		if err != nil {
			return err
		}
		task.UserBoardRoleID = &ubrObj.ID
	}

	// check permission for creator
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, user.ID, board.ID)
	if err != nil {
		return err
	}

	if !rbac.HasPermission(role, rbac.PermissionCreateTask) {
		return ErrPermissionDenied
	}

	col, err := s.columnOps.GetMinOrderColumn(ctx, task.BoardID)
	if err != nil {
		return err
	}
	task.ColumnID = col.ID
	err = s.taskOps.Create(ctx, task)
	if err != nil {
		return err
	}

	// notif to owner and maintainer!!! TO Do
	return nil
}

func (s *TaskService) AddDependency(ctx context.Context, task *t.Task) error {
	// task exists?
	existedTask, err := s.taskOps.GetTaskByID(ctx, task.ID)
	if err != nil {
		return err
	}
	// check permission
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, task.CreatedByUserID, existedTask.BoardID)
	if err != nil {
		return ErrPermissionDenied
	}

	if !rbac.HasPermission(role, rbac.PermissionCreateTask) {
		return ErrPermissionDenied
	}

	return s.taskOps.AddDependency(ctx, task)
}

func (s *TaskService) GetFullTaskByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*t.Task, error) {
	task, err := s.taskOps.GetFullTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	fetcherRole, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, task.BoardID)
	if err != nil {
		return nil, ErrPermissionDenied
	}

	if !rbac.HasPermission(fetcherRole, rbac.PermissionViewTask) {
		return nil, ErrPermissionDenied
	}

	return task, err
}

func (s *TaskService) UpdateTaskColumnByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, colID uuid.UUID) (*t.Task, error) {
	task, err := s.taskOps.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	fetcherRole, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, task.BoardID)
	if err != nil {
		return nil, ErrPermissionDenied
	}

	if !rbac.HasPermission(fetcherRole, rbac.PermissionMoveOwnTask) {
		return nil, ErrPermissionDenied
	}

	updatedTask, err := s.taskOps.UpdateTaskColumnByID(ctx, taskID, colID)
	if err != nil {
		return nil, err
	}

	b, err := s.boardOps.GetBoardByID(ctx, task.BoardID)
	if err != nil {
		return nil, err
	}
	updater, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	newColumn, err := s.columnOps.GetColumnByID(ctx, updatedTask.ColumnID)
	if err != nil {
		return nil, err
	}
	userBoardRoleObj, err := s.userBoardRoleOps.GetUserBoardRoleObj(ctx, userID, task.BoardID)
	if err != nil {
		return nil, err
	}
	description := fmt.Sprintf("Task %s from Board %s Moved to Column %s By %s", task.Title, b.Name, newColumn.Name, updater.FirstName)

	newNotification := notification.NewNotification(description, notification.TaskMoved, userBoardRoleObj.ID)

	err = s.notificaionOps.NotifBroadCasting(ctx, newNotification, task.BoardID, userID, task)
	if err != nil {
		return nil, err
	}
	return updatedTask, err
}

func (s *TaskService) ReorderTasks(ctx context.Context, userID, colID uuid.UUID, newOrder map[uuid.UUID]uint) ([]t.Task, error) {
	col, err := s.columnOps.GetColumnByID(ctx, colID)
	if err != nil {
		return nil, column.ErrColumnNotFound
	}
	role, err := s.userBoardRoleOps.GetUserBoardRole(ctx, userID, col.BoardID)
	if err != nil {
		return nil, ErrPermissionDenied
	}

	if !rbac.HasPermission(role, rbac.PermissionViewTask) {
		return nil, ErrPermissionDenied
	}
	tasks, err := s.taskOps.ReorderTasks(ctx, colID, newOrder)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
