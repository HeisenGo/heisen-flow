package handlers

import (
	"errors"
	"fmt"
	presenter "server/api/http/handlers/presentor"
	"server/internal/user"
	"server/service"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.UserRegisterReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		u := presenter.UserRegisterToUserDomain(&req)

		new_user, err := authService.CreateUser(c.Context(), u)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, user.ErrInvalidEmail) || errors.Is(err, user.ErrInvalidPassword) || errors.Is(err, user.ErrEmailAlreadyExists) {
				status = fiber.StatusBadRequest
			}

			return SendError(c, err, status)
		}

		return Created(c, "user successfully registered", fiber.Map{
			"user_id": new_user.ID,
		})
	}
}

func LoginUser(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		c.Cookie(&fiber.Cookie{
			Name:        "X-Session-ID",
			Value:       fmt.Sprint(time.Now().UnixNano()),
			HTTPOnly:    true,
			SessionOnly: true,
		})

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		authToken, err := authService.Login(c.Context(), input.Email, input.Password)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		return SendUserToken(c, authToken)
	}
}

func RefreshToken(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refToken := c.GetReqHeaders()["Authorization"][0]
		if len(refToken) == 0 {
			return SendError(c, errors.New("token should be provided"), fiber.StatusBadRequest)
		}
		pureToken := strings.Split(refToken, " ")[1]
		authToken, err := authService.RefreshAuth(c.UserContext(), pureToken)
		if err != nil {
			return SendError(c, err, fiber.StatusUnauthorized)
		}

		return SendUserToken(c, authToken)
	}
}
