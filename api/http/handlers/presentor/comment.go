package presenter

import (
	"server/internal/comment"
	userboardrole "server/internal/user_board_role"

	"github.com/google/uuid"
)

type CommentCreateReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description" `
	TaskID      uuid.UUID `json:"task_id"`
}

func CommentReqToCommentDomain(up *CommentCreateReq, userID uuid.UUID) (*userboardrole.UserBoardRole, *comment.Comment) {

	ubr := &userboardrole.UserBoardRole{
		UserID: userID,
	}

	c := &comment.Comment{
		Title:       up.Title,
		Description: up.Description,
		TaskID:      up.TaskID,
	}

	return ubr, c
}
