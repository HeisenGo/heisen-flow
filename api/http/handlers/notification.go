package handlers

import (
	"errors"
	presenter "server/api/http/handlers/presentor"
	"server/internal/notification"
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
			return presenter.BadRequest(c, errWrongClaimType)
		}
		notifList, err := notificationService.GetUserNotifications(c.UserContext(), userClaims.UserID)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, notification.ErrNotifsNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.BatchNotifToNotifResp(notifList)
		return presenter.OK(c, "notifications successfully fetched", res)
	}
}

func UpdateNotifications(notificationService *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return presenter.BadRequest(c, errWrongClaimType)
		}

		notificationIDParam := c.Params("notifID")
		notificationID, err := uuid.Parse(notificationIDParam)
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		n, err := notificationService.MarkNotificationAsSeen(c.UserContext(), notificationID, userClaims.UserID)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, notification.ErrNotifNotFound) || errors.Is(err, service.ErrPermissionDenied) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.DomainNotifToNotifResp(*n)
		return presenter.OK(c, "Marked As Seen", res)
	}
}
