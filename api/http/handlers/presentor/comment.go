package presenter

import (
	"server/internal/comment"
	"time"

	"github.com/google/uuid"
)

type CommentCreateReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TaskID      uuid.UUID `json:"task_id"`
}

func CommentReqToCommentDomain(up *CommentCreateReq) *comment.Comment {
	return &comment.Comment{
		Title:       up.Title,
		Description: up.Description,
		TaskID:      up.TaskID,
	}
}

type CommentCreateRep struct {
	ID          uuid.UUID `json:"comment_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TaskID      uuid.UUID `json:"task_id"`
	CreatedAt   time.Time   `json:"created_at"`
}

func CommentToCommentCreateResp(c *comment.Comment) *CommentCreateRep {
	return &CommentCreateRep{
		ID:          c.ID,
		CreatedAt:   c.CreatedAt,
		Title:       c.Title,
		Description: c.Description,
		TaskID:      c.TaskID,
	}
}
