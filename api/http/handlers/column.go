package handlers

import (
	presenter "server/api/http/handlers/presentor"
	"server/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateColumns(serviceFactory ServiceFactory[*service.ColumnService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		columnService := serviceFactory(c.UserContext())

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
		createdColumns, err := columnService.CreateColumns(c.UserContext(), columns)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.EntitiesToCreateColumnsResponse(createdColumns)
		return presenter.Created(c, resp.Message, fiber.Map{"data": resp.Data})
	}
}

func DeleteColumn(serviceFactory ServiceFactory[*service.ColumnService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		columnService := serviceFactory(c.UserContext())

		columnIDParam := c.Params("columnID")
		columnID, err := uuid.Parse(columnIDParam)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = columnService.DeleteColumn(c.UserContext(), columnID)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Column deleted successfully", nil)
	}
}

// func GetColumnsByBoardID(serviceFactory ServiceFactory[*service.ColumnService]) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		columnService := serviceFactory(c.UserContext())

// 		boardIDParam := c.Params("boardID")
// 		boardID, err := uuid.Parse(boardIDParam)
// 		if err != nil {
// 			return SendError(c, err, fiber.StatusBadRequest)
// 		}

// 		columns, err := columnService.GetColumnsByBoardID(c.UserContext(), boardID)
// 		if err != nil {
// 			return InternalServerError(c, err)
// 		}

// 		resp := presenter.EntitiesToGetColumnsResponse(columns)
// 		return OK(c, resp.Message, fiber.Map{"data": resp.Columns})
// 	}
// }
