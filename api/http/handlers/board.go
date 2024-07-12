package handlers

import (
	"errors"
	"github.com/google/uuid"
	presenter "server/api/http/handlers/presentor"
	"server/internal/board"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetUserBoards(boardService *service.BoardService) fiber.Handler {
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
		data := presenter.NewPagination(
			presenter.BatchBoardsToUserBoard(boards),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "boards successfully fetched.", data)
	}
}

func GetPublicBoards(boardService *service.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		//query parameter
		page, pageSize := PageAndPageSize(c)

		boards, total, err := boardService.GetPublicBoards(c.UserContext(), userClaims.UserID, uint(page), uint(pageSize))
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		data := presenter.NewPagination(
			presenter.BatchBoardsToUserBoard(boards),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "boards successfully fetched.", data)
	}
}

func GetFullBoardByID(boardService *service.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		boardID, err := uuid.Parse(c.Params("boardID"))
		if err != nil {
			return presenter.BadRequest(c, errors.New("given board_id format in path is not correct"))
		}
		fetchedBoard, err := boardService.GetFullBoardByID(c.UserContext(), userClaims.UserID, boardID)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		data := presenter.BoardToFullBoardResp(*fetchedBoard)
		return presenter.OK(c, "board successfully fetched.", data)
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
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "Board created successfully", fiber.Map{
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
			if errors.Is(err, service.ErrPermissionDeniedToInvite) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, service.ErrAMember) || errors.Is(err, service.ErrOwnerExists) || errors.Is(err, service.ErrUndefinedRole) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "User successfully invited", fiber.Map{
			"role":               ubr.Role,
			"user_board_role_id": ubr.ID,
			"board_id":           ubr.BoardID,
			"email":              req.Email,
		})
	}
}
