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
// GetNotifications retrieves notifications for the authenticated user.
// @Summary Get user notifications
// @Description Retrieve a list of notifications for the authenticated user.
// @Tags Notifications
// @Accept  json
// @Produce  json
// @Success 200 {object} presenter.NotifResp "Notifications successfully fetched"
// @Failure 400 {object} map[string]interface{} "Bad request, invalid user claims or user not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /notifications [get]
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
// UpdateNotifications marks a notification as seen for the authenticated user.
// @Summary Mark notification as seen
// @Description Marks a specific notification as seen for the authenticated user.
// @Tags Notifications
// @Accept  json
// @Produce  json
// @Param notifID path string true "Notification ID"
// @Success 200 {object} presenter.NotifResp "Notification marked as seen"
// @Failure 400 {object} map[string]interface{} "Bad request, invalid user claims, or notification ID format"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /notifications/{notifID} [put]
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
