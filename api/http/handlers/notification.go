package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

func GetNotifications(notificationService *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		notifList , err := notificationService.GetUserUnseenNotifications(c.UserContext(),userClaims.UserID )
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		return c.JSON(fiber.Map{"notifications": notifList})
	}
}

func UpdateNotifications(notificationService *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.NotificationReq
		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		notificationIDParam := c.Params("notificationID")
		notificationID, err := uuid.Parse(notificationIDParam)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		err = notificationService.MarkNotificationAsSeen(c.UserContext(),notificationID)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		return presenter.OK(c,"Marked As Seen",err)
	}
}