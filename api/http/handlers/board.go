package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/board"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

// GetUserBoards retrieves the boards of the logged-in user.
// @Summary Get user's boards
// @Description Retrieve the boards associated with the authenticated user.
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.BoardUserResp "boards: paginated list of user's boards"
// @Failure 400 {object} map[string]interface{} "error: bad request, wrong claim type"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /user/boards [get]
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

// GetPublicBoards retrieves public boards.
// @Summary Get public boards
// @Description Retrieve the public boards.
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{} "boards: paginated list of public boards"
// @Failure 400 {object} map[string]interface{} "error: bad request, wrong claim type"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /public/boards [get]
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

// GetFullBoardByID retrieves a full board by its ID.
// @Summary Get full board by ID
// @Description Retrieve a full board by its ID for the authenticated user.
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param boardID path string true "Board ID"
// @Success 200 {object} presenter.FullBoardResp "board: the full board details"
// @Failure 400 {object} map[string]interface{} "error: bad request, wrong claim type or invalid board ID format"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /boards/{boardID} [get]
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

// CreateUserBoard creates a new board for the user.
// @Summary Create user board
// @Description Create a new board for the authenticated user.
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param board body presenter.UserBoard true "Board details"
// @Success 201 {object} presenter.CreateBoardResponse "board: the created board details"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid board details"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /user/boards [post]
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
			if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, board.ErrWrongType) || errors.Is(err, board.ErrInvalidName) {
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}
		res := presenter.BoardToCreateBoardResponse(b)
		return presenter.Created(c, "Board created successfully", res)
	}
}

// InviteToBoard invites a user to a board.
// @Summary Invite user to board
// @Description Invite a user to a board with a specified role.
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param invite body presenter.InviteUserToBoard true "Invitation details"
// @Success 200 {object} presenter.InviteMemberResponse "invite: the details of the invitation"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid invitation details"
// @Failure 403 {object} map[string]interface{} "error: forbidden, permission denied to invite"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /boards/invite [post]
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
			if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, service.ErrAMember) || errors.Is(err, service.ErrOwnerExists) || errors.Is(err, service.ErrUndefinedRole) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.InviteMemberToInviteMemberResponse(ubr, req.Email)

		return presenter.Created(c, "User successfully invited", res)
	}
}

func DeleteBoard(boardService *service.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		boardIDParam := c.Params("boardID")
		boardID, err := uuid.Parse(boardIDParam)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		ubr := presenter.DeleteBoardParamToUserBoardRole(boardID, userClaims.UserID)
		err = boardService.DeleteBoardByID(c.UserContext(), ubr)

		if err != nil {
			if errors.Is(err, service.ErrPermissionDeniedToDelete) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, board.ErrBoardNotFound) || errors.Is(err, user.ErrUserNotFound) {
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}

		return presenter.NoContent(c)
	}
}
