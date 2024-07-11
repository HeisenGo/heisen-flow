package handlers

import (
	presenter "server/api/http/handlers/presentor"
	"server/service"
	"github.com/gofiber/fiber/v2"
)

func CreateNotification(notificationService *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.NotificationReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		
		notif := presenter.NotificationToNotificationDomain(&req)
		if err := notificationService.CreateNotification(c.UserContext(), notif); err != nil {
				return BadRequest(c, err)
			}
		
		return Created(c, "Notification created successfully", fiber.Map{
			"notification_id":        notif.ID,
		})
	}
}