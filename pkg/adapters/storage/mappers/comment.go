package mappers

import (
	"server/internal/comment"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"
)

func CommentEntityToDomain(commentEntity entities.Comment) comment.Comment {
	return comment.Comment{
		Title:           commentEntity.Title,
		Description:     commentEntity.Description,
		UserBoardRoleID: commentEntity.UserBoardRoleID,
		TaskID:          commentEntity.TaskID,
		CreatedAt:       commentEntity.CreatedAt,
	}
}

func BatchCommentEntitiesToDomain(commentEntities []entities.Comment) []comment.Comment {
	return fp.Map(commentEntities, CommentEntityToDomain)
}

func CommentDomainToEntity(c *comment.Comment) *entities.Comment {
	return &entities.Comment{
		Title:           c.Title,
		Description:     c.Description,
		UserBoardRoleID: c.UserBoardRoleID,
		TaskID:          c.TaskID,
	}
}
