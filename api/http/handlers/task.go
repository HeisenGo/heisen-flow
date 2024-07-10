package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/board"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"

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

func CreateTask(serviceFactory ServiceFactory[*service.TaskService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskService := serviceFactory(c.UserContext())

		var req presenter.UserTask

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		t := presenter.UserTaskToTask(&req, userClaims.UserID)

		if err := taskService.CreateTask(c.UserContext(), t); err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, service.ErrPermissionDenied) {
				status = fiber.StatusForbidden
			}
			if errors.Is(err, service.ErrNotMember) || errors.Is(err, user.ErrUserNotFound) || errors.Is(err, board.ErrBoardNotFound) {
				status = fiber.StatusBadGateway
			}

			return SendError(c, err, status)
		}

		return c.JSON(fiber.Map{
			"message": "task created",
			"task_id": t.ID,
		})
	}
}
