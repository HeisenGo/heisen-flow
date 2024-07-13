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

// CreateColumns creates multiple columns for a board.
// @Summary Create columns
// @Description Create multiple columns for a specified board.
// @Tags Columns
// @Accept  json
// @Produce  json
// @Param board_id body string true "Board ID"
// @Param columns body presenter.CreateColumnsRequest true "Columns creation details"
// @Success 201 {object} presenter.CreateColumnsResponse "response: details of created columns"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid board ID format or missing columns details"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Security BearerAuth
// @Router /columns [post]
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
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.EntitiesToCreateColumnsResponse(createdColumns)
		return presenter.Created(c, resp.Message, resp.Data)
	}
}

// DeleteColumn deletes a column by its ID.
// @Summary Delete column
// @Description Delete a column by its ID.
// @Tags Columns
// @Accept  json
// @Produce  json
// @Param columnID path string true "Column ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid column ID format"
// @Failure 404 {object} map[string]interface{} "error: not found, column not found"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Security BearerAuth
// @Router /columns/{columnID} [delete]
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

// ReorderColumns reorders the columns of a board.
// @Summary Reorder columns
// @Description Reorder the columns of a board for the authenticated user.
// @Tags Columns
// @Accept  json
// @Produce  json
// @Param ReorderColumnsRequest body presenter.ReorderColumnsRequest true "Reorder Columns Request"
// @Success 200 {object} []presenter.ColumnResponseItem "Columns reordered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request, invalid reorder details"
// @Failure 403 {object} map[string]interface{} "Forbidden, permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /columns [put]
func ReorderColumns(serviceFactory ServiceFactory[*service.ColumnService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		columnService := serviceFactory(c.UserContext())

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var req presenter.ReorderColumnsRequest
		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		boardID, newOrder := presenter.ReorderColumnsRequestToMap(req)
		cols, err := columnService.ReorderColumns(c.UserContext(), userClaims.UserID, boardID, newOrder)
		if err != nil {
			if errors.Is(err, service.ErrPermissionDeniedToDeleteColumn) {
				presenter.Forbidden(c, err)
			}
			if errors.Is(err, column.ErrColumnNotFound) || errors.Is(err, column.ErrFailedToFetchColumns) || errors.Is(err, column.ErrFailedToUpdateColumn) || errors.Is(err, column.ErrInvalidColumnID) || errors.Is(err, column.ErrLengthMismatch) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.BatchColumnToColumnResponseItem(cols)
		return presenter.OK(c, "Columns ReOrdered Successfully", res)
	}
}
