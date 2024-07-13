package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/user"
	"server/pkg/jwt"
	"server/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetNotifications(notificationService *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		notifList, err := notificationService.GetUserNotifications(c.UserContext(), userClaims.UserID)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		res := presenter.BatchNotifToNotifResp(notifList)
		return presenter.OK(c, "notifications successfully fetched", res)
	}
}

func UpdateNotifications(notificationService *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		
		notificationIDParam := c.Params("notifID")
		notificationID, err := uuid.Parse(notificationIDParam)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		n, err := notificationService.MarkNotificationAsSeen(c.UserContext(), notificationID, userClaims.UserID)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrUserNotFound) {
				status = fiber.StatusBadRequest
			}
			return SendError(c, err, status)
		}
		res := presenter.DomainNotifToNotifResp(*n)
		return presenter.OK(c, "Marked As Seen", res)
	}
}
