package service

import (
	"context"
	"server/internal/comment"
	"server/internal/notification"
	t "server/internal/task"
	userboardrole "server/internal/user_board_role"
	"server/pkg/rbac"

	"github.com/google/uuid"
)

// CommentService handles board-related operations

type CommentService struct {
	commentOps       *comment.Ops
	userBoardRoleOps *userboardrole.Ops
	notifOps         *notification.Ops
	taskOps          *t.Ops
}

// NewCommentService creates a new BoardService

func NewCommentService(commentOps *comment.Ops, userBoardOps *userboardrole.Ops, notifOps *notification.Ops, taskOps *t.Ops) *CommentService {
	return &CommentService{
		commentOps:       commentOps,
		userBoardRoleOps: userBoardOps,
		notifOps:         notifOps,
		taskOps:          taskOps,
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

	return nil
}
