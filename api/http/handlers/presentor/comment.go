package presenter

import (
	"github.com/google/uuid"
	"server/internal/comment"
	userBoardRole "server/internal/user_board_role"

)

type CommentCreateReq struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	TaskID          uuid.UUID `json:"task_id"`
}

func CommentReqToCommentDomain(up *CommentCreateReq, id uuid.UUID) (*comment.Comment,*userBoardRole.UserBoardRole) {
	c:= &comment.Comment{
		Title:           up.Title,
		Description:     up.Description,
		TaskID:          up.TaskID,
	}

	ubr := &userBoardRole.UserBoardRole{
		UserID: id,
	}

	return c, ubr
}
