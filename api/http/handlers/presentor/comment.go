package presenter

import (
	"github.com/google/uuid"
	"server/internal/comment"
)

type CommentCreateReq struct {
	Title           string    `json:"title"`
	Description     string    `json:"description" `
	UserBoardRoleID uuid.UUID `json:"user_board_role_id"`
	TaskID          uuid.UUID `json:"task_id"`
	//CreatedAt       Timestamp
}

func CommentReqToCommentDomain(up *CommentCreateReq) *comment.Comment {
	return &comment.Comment{
		Title:           up.Title,
		Description:     up.Description,
		UserBoardRoleID: up.UserBoardRoleID,
		TaskID:          up.TaskID,
	}
}
