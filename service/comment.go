package service

import (
	"context"
	"server/internal/comment"
	userboardrole "server/internal/user_board_role"
)

// CommentService handles board-related operations

type CommentService struct {
	commentOps       *comment.Ops
	userBoardRoleOps *userboardrole.Ops
}

// NewCommentService creates a new BoardService

func NewCommentService(commentOps *comment.Ops, userBoardOps *userboardrole.Ops) *CommentService {
	return &CommentService{commentOps: commentOps, userBoardRoleOps: userBoardOps}
}

func (s *CommentService) CreateComment(ctx context.Context, c *comment.Comment, ubr *userboardrole.UserBoardRole) error {

	//Validate Comment

	err := s.commentOps.Insert(ctx, c)
	if err != nil {
		return err
	}

	return nil
}
