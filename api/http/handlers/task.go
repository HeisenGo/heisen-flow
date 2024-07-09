package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

// func UserTasksInABoard(orderService *service.OrderService) fiber.Handler {
//
// }

// all tasks of a board till first depth!

// func BoardTasks(orderService *service.OrderService) fiber.Handler {
//
// }

// func TasksSubTasks(orderService *service.OrderService) fiber.Handler {
//
// }

func CreateTask(serviceFactory ServiceFactory[*service.Task]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orderService := serviceFactory(c.UserContext())

		var req presenter.UserTask

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		t := presenter.UserTaskToTask(&req, userClaims.UserID)

		if err := orderService.CreateOrder(c.UserContext(), t); err != nil {
			status := fiber.StatusInternalServerError
			// if errors.Is(err, .ErrQuantityGreater) || errors.Is(err, order.ErrWrongOrderTime) {
			// 	status = fiber.StatusBadRequest
			// }

			return SendError(c, err, status)
		}

		return c.JSON(fiber.Map{
			"task_id": t.ID,
		})
	}
}
