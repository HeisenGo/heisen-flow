package handlers

import (
	"github.com/gofiber/fiber/v2"
	presenter "server/api/http/handlers/presentor"
	"server/pkg/jwt"
	"server/service"
)

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

		if err := commentService.CreateComment(c.UserContext(), cr, userClaims.UserID); err != nil { //here
			//if errors.Is(err, comment.ErrWrongType) || errors.Is(err, comment.ErrInvalidName) {
			//	return BadRequest(c, err)
			//}

			return InternalServerError(c, err)
		}

		return Created(c, "Comment created successfully", fiber.Map{})
	}
}
