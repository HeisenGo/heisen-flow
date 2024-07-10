package handlers

import (
	presenter "server/api/http/handlers/presentor"
	"server/service"

	"github.com/gofiber/fiber/v2"
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
			return InternalServerError(c, err)
		}

		columns := presenter.CreateColumnsRequestToEntities(req, maxOrder)
		createdColumns, err := columnService.CreateColumns(c.UserContext(), columns)
		if err != nil {
			return InternalServerError(c, err)
		}

		resp := presenter.EntitiesToCreateColumnsResponse(createdColumns)
		return Created(c, resp.Message, fiber.Map{"data": resp.Data})
	}
}
