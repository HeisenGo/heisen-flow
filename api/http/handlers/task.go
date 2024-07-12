package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/board"
	"server/internal/task"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTask(serviceFactory ServiceFactory[*service.TaskService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskService := serviceFactory(c.UserContext())

		var req presenter.UserTask

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
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
			if errors.Is(err, service.ErrNotMember) || errors.Is(err, user.ErrUserNotFound) || errors.Is(err, board.ErrBoardNotFound) || errors.Is(err, service.ErrCantAssigned) || errors.Is(err, task.ErrInvalidStoryPoint) {
				status = fiber.StatusBadGateway
			}

			return SendError(c, err, status)
		}

		return presenter.Created(c, "Task created successfully", fiber.Map{
			"message": "task created",
			"task_id": t.ID,
		})
	}
}

func AddDependency(serviceFactory ServiceFactory[*service.TaskService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskService := serviceFactory(c.UserContext())

		var req presenter.DependentTasks

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}

		t := presenter.AddDependencyReqToTask(&req, userClaims.UserID)

		if err := taskService.AddDependency(c.UserContext(), t); err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, service.ErrPermissionDenied) {
				status = fiber.StatusForbidden
			}
			if errors.Is(err, task.ErrCircularDependency) || errors.Is(err, task.ErrTaskNotFound) || errors.Is(err, task.ErrFailedToFindDependsOnTasks) {
				status = fiber.StatusBadGateway
			}

			return SendError(c, err, status)
		}

		return presenter.Created(c, "Dependency added successfully", fiber.Map{
			"message": "task dependency created",
			"task_id": t.ID,
		})
	}
}
