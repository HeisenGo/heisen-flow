package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/task"
	"server/pkg/jwt"
	"server/service"

	"github.com/gofiber/fiber/v2"
)
// CreateUserComment creates a new comment on a task.
// @Summary Create comment
// @Description Creates a new comment on a specific task for the authenticated user.
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body presenter.CommentCreateReq true "Comment Create Request"
// @Success 201 {object} presenter.CommentCreateRep "Comment created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request, invalid user claims, or task not found"
// @Failure 403 {object} map[string]interface{} "Forbidden, user does not have permission to create a comment"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /comments [post]
func CreateUserComment(serviceFactory ServiceFactory[*service.CommentService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		commentService := serviceFactory(c.UserContext())

		var req presenter.CommentCreateReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		cr := presenter.CommentReqToCommentDomain(&req)

		err := commentService.CreateComment(c.UserContext(), cr, userClaims.UserID)
		if err != nil {
			if errors.Is(err, service.ErrPermissionDenied) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, task.ErrTaskNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		resp := presenter.CommentToCommentCreateResp(cr)
		return presenter.Created(c, "Comment created successfully", resp)
	}
}
