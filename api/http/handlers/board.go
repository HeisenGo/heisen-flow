package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/board"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UserBoards(boardService *service.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		//query parameter
		page, pageSize := PageAndPageSize(c)

		boards, total, err := boardService.GetUserBoards(c.UserContext(), userClaims.UserID, uint(page), uint(pageSize))
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}

		return c.JSON(fiber.Map{"boards": boards, "total": total}) // needs a presenter !!!! To Do
	}
}

func CreateUserBoard(serviceFactory ServiceFactory[*service.BoardService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardService := serviceFactory(c.UserContext())

		var req presenter.UserBoard

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		b, ubr := presenter.UserBoardToBoard(&req, userClaims.UserID)
		b.CreatedAt = time.Now()
		if err := boardService.CreateBoard(c.UserContext(), b, ubr); err != nil {
			if errors.Is(err, board.ErrWrongType) || errors.Is(err, board.ErrInvalidName) {
				return BadRequest(c, err)
			}

			return InternalServerError(c, err)
		}

		return Created(c, "Board created successfully", fiber.Map{
			"board_id":        b.ID,
			"user_board_role": ubr.ID,
		})
	}
}

func InviteToBoard(serviceFactory ServiceFactory[*service.BoardService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardService := serviceFactory(c.UserContext())

		var req presenter.InviteUserToBoard

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		ubr := presenter.InviteUserToBoardToUserBoardRole(&req)
		if err := boardService.InviteUser(c.UserContext(), userClaims.UserID, req.Email, ubr); err != nil {
			//status := fiber.StatusInternalServerError
			/// needs to be checked by error wrapping
			// if errors.Is(err, board.ErrWrongType) || errors.Is(err, board.ErrInvalidName) {
			// 	status = fiber.StatusBadRequest
			// }

			return BadRequest(c, err)
		}

		return OK(c, "User successfully invited", fiber.Map{
			"role":               ubr.Role,
			"user_board_role_id": ubr.ID,
			"board_id":           ubr.BoardID,
			"email":              req.Email,
		})
	}
}
