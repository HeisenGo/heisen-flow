package service

import (
	"context"
	"fmt"
	"server/internal/board"
	"server/internal/comment"
	"server/internal/notification"
	t "server/internal/task"
	"server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

// CommentService handles board-related operations

type CommentService struct {
	userOps          *user.Ops
	commentOps       *comment.Ops
	userBoardRoleOps *userboardrole.Ops
	notifOps         *notification.Ops
	taskOps          *t.Ops
	boardOps         *board.Ops
}

// NewCommentService creates a new BoardService

func NewCommentService(commentOps *comment.Ops, userBoardOps *userboardrole.Ops, notifOps *notification.Ops,
	taskOps *t.Ops, userOps *user.Ops, boardOps *board.Ops) *CommentService {
	return &CommentService{
		commentOps:       commentOps,
		userBoardRoleOps: userBoardOps,
		notifOps:         notifOps,
		taskOps:          taskOps,
		userOps:          userOps,
		boardOps:         boardOps,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, c *comment.Comment, userID uuid.UUID) error {
	task, err := s.taskOps.GetTaskByID(ctx, c.TaskID)
	if err != nil {
		return t.ErrTaskNotFound
	}
	userBoardRoleObj, err := s.userBoardRoleOps.GetUserBoardRoleObj(ctx, userID, task.BoardID)
	if err != nil {
		return ErrPermissionDenied
	}
	if userBoardRoleObj.ID != *task.UserBoardRoleID && userBoardRoleObj.Role == string(rbac.RoleEditor) {
		return ErrPermissionDenied
	}
	if !rbac.HasPermission(rbac.Role(userBoardRoleObj.Role), rbac.PermissionCommentOwnTask) || !rbac.HasPermission(rbac.Role(userBoardRoleObj.Role), rbac.PermissionCommentAnyTask) {
		return ErrPermissionDenied
	}
	c.UserBoardRoleID = userBoardRoleObj.ID
	err = s.commentOps.Insert(ctx, c)
	if err != nil {
		return err
	}
	// send notif to maintainaers owners and asignee of task
	// Assignee : userBoardRoleObj.UserID
	commenter, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	board, err := s.boardOps.GetBoardByID(ctx, task.BoardID)
	if err != nil {
		return err
	}
	description := fmt.Sprintf("%s commented on task '%s' of board '%s'", commenter.FirstName, task.Title, board.Name)
	notif := notification.NewNotification(description, notification.CommentedNotif, userBoardRoleObj.ID)
	// TODO : editor comment notif
	err = s.notifOps.NotifBroadCasting(ctx, notif, board.ID, userID, task)
	if err != nil {
		return err
	}
	return nil
}
