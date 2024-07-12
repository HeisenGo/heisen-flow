package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/board"
	"server/internal/column"
	"server/pkg/jwt"
	"server/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateColumns(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var req presenter.CreateColumnsRequest
		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		// Get the maximum order of existing columns for the board
		maxOrder, err := columnService.GetMaxOrderForBoard(c.UserContext(), req.BoardID)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		columns := presenter.CreateColumnsRequestToEntities(req, maxOrder)
		createdColumns, err := columnService.CreateColumns(c.UserContext(), columns, userClaims.UserID)
		if err != nil {
			if errors.Is(err, service.ErrPermissionDeniedToDeleteColumn) {
				presenter.Forbidden(c, err)
			}
			if errors.Is(err, board.ErrBoardNotFound) {
				presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.EntitiesToCreateColumnsResponse(createdColumns)
		return presenter.Created(c, resp.Message, resp.Data)
	}
}

func DeleteColumn(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		columnIDParam := c.Params("columnID")
		columnID, err := uuid.Parse(columnIDParam)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = columnService.DeleteColumn(c.UserContext(), columnID, userClaims.UserID)
		if err != nil {
			if errors.Is(err, service.ErrPermissionDeniedToDeleteColumn) {
				presenter.Forbidden(c, err)
			}
			if errors.Is(err, column.ErrColumnNotEmpty) {
				return presenter.BadRequest(c, err)
			}
			if errors.Is(err, column.ErrColumnNotFound) {
				return presenter.NotFound(c, err)
			}
			return presenter.InternalServerError(c, err)
		}

		return presenter.NoContent(c)
	}
}
