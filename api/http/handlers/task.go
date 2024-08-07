package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/board"
	"server/internal/column"
	"server/internal/task"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

// CreateTask creates a new task.
// @Summary Create task
// @Description Create a new task for the authenticated user.
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param task body presenter.UserTask true "Task details"
// @Success 201 {object} presenter.CreateTaskResp "response: details of created task"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid task details"
// @Failure 403 {object} map[string]interface{} "error: forbidden, permission denied"
// @Failure 502 {object} map[string]interface{} "error: bad gateway, not a member, user not found, board not found, or other error"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Security BearerAuth
// @Router /tasks [post]
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
		res := presenter.DomainTaskToCreateTaskResp(t)
		return presenter.Created(c, "Task created successfully", res)
	}
}

// AddDependency adds a dependency between tasks.
// @Summary Add task dependency
// @Description Add a dependency between tasks for the authenticated user.
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param dependency body presenter.DependentTasks true "Dependency details"
// @Success 201 {object} map[string]interface{} "response: details of added task dependency"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid dependency details"
// @Failure 403 {object} map[string]interface{} "error: forbidden, permission denied"
// @Failure 502 {object} map[string]interface{} "error: bad gateway, circular dependency, task not found, or other error"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Security BearerAuth
// @Router /tasks/dependency [post]
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

// GetFullTaskByID retrieves the full details of a task by its ID.
// @Summary Get full task by ID
// @Description Retrieve the full details of a task by its ID for the authenticated user.
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param taskID path string true "Task ID"
// @Success 200 {object} presenter.FullTaskResp "Task successfully fetched"
// @Failure 400 {object} map[string]interface{} "Bad request, invalid task ID or user not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /tasks/{taskID} [get]
func GetFullTaskByID(taskService *service.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		taskID, err := uuid.Parse(c.Params("taskID"))
		if err != nil {
			return presenter.BadRequest(c, errors.New("given task_id format in path is not correct"))
		}
		fetchedTask, err := taskService.GetFullTaskByID(c.UserContext(), userClaims.UserID, taskID)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		data := presenter.TaskToFullTaskResp(*fetchedTask)
		return presenter.OK(c, "task successfully fetched.", data)
	}
}

// UpdateTaskColumnByID updates the column of a task by its ID.
// @Summary Update task column by ID
// @Description Update the column of a task by its ID for the authenticated user.
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param taskID path string true "Task ID"
// @Param UpdateTaskColReq body presenter.UpdateTaskColReq true "Update Task Column Request"
// @Success 200 {object} presenter.UpdatedTaskResp "Task successfully updated"
// @Failure 400 {object} map[string]interface{} "Bad request, invalid task ID or column ID, or dependent task issues"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /tasks/{taskID}/column [put]
func UpdateTaskColumnByID(serviceFactory ServiceFactory[*service.TaskService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskService := serviceFactory(c.UserContext())
		var req presenter.UpdateTaskColReq

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
		taskID, err := uuid.Parse(c.Params("taskID"))
		if err != nil {
			return presenter.BadRequest(c, errors.New("given task_id format in path is not correct"))
		}
		updatedTask, err := taskService.UpdateTaskColumnByID(c.UserContext(), userClaims.UserID, taskID, req.ColumnID)
		if err != nil {
			if errors.Is(err, task.ErrTaskNotFound) || errors.Is(err, task.ErrColumnNotFound) || errors.Is(err, task.ErrCantDoneDependentTask) {
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}
		data := presenter.TaskToUpdatedTaskResp(*updatedTask)
		return presenter.OK(c, "task successfully updated.", data)
	}
}

// ReorderTasks reorders the tasks of a board.
// @Summary Reorder Tasks
// @Description Reorder the tasks of a board for the authenticated user.
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param ReorderColumnsRequest body presenter.ReorderTasksReq true "Reorder Tasks Request"
// @Success 200 {object} []presenter.TaskReorderRespItem "Tasks reordered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request, invalid reorder details"
// @Failure 403 {object} map[string]interface{} "Forbidden, permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /tasks/reorder [patch]
func ReorderTasks(serviceFactory ServiceFactory[*service.TaskService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskService := serviceFactory(c.UserContext())

		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		var req presenter.ReorderTasksReq
		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		colID, newOrder := presenter.ReorderTasksReqToMap(req)
		tasks, err := taskService.ReorderTasks(c.UserContext(), userClaims.UserID, colID, newOrder)
		if err != nil {
			if errors.Is(err, service.ErrPermissionDeniedToDeleteColumn) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, column.ErrColumnNotFound) || errors.Is(err, column.ErrFailedToFetchColumns) || errors.Is(err, column.ErrFailedToUpdateColumn) || errors.Is(err, column.ErrInvalidColumnID) || errors.Is(err, column.ErrLengthMismatch) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.BatchTaskToTaskReorderRespItem(tasks)
		return presenter.OK(c, "Tasks ReOrdered Successfully", res)
	}
}
