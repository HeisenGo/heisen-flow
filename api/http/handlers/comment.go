package handlers

import (
	presenter "server/api/http/handlers/presentor"
	"server/pkg/jwt"
	"server/service"

	"github.com/gofiber/fiber/v2"
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

		ubr, cr := presenter.CommentReqToCommentDomain(&req, userClaims.UserID)

		if err := commentService.CreateComment(c.UserContext(), cr, ubr); err != nil { //here
			//if errors.Is(err, comment.ErrWrongType) || errors.Is(err, comment.ErrInvalidName) {
			//	return BadRequest(c, err)
			//}

			return InternalServerError(c, err)
		}

		return Created(c, "Comment created successfully", fiber.Map{})
	}
}
