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
// RegisterUser registers a new user.
// @Summary Register a new user
// @Description Create a new user account with the provided details.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body presenter.UserRegisterReq true "User registration details"
// @Success 201 {object} map[string]interface{} "user_id: the ID of the newly registered user"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid email or password, or email already exists"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /register [post]
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
// LoginUser logs in an existing user.
// @Summary Login an existing user
// @Description Authenticate a user with email and password.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param login body object true "Login details"
// @Success 200 {object} map[string]interface{} "auth_token: the authentication token for the user"
// @Failure 400 {object} map[string]interface{} "error: bad request, invalid email or password"
// @Router /login [post]
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
// RefreshToken refreshes the authentication token.
// @Summary Refresh authentication token
// @Description Refresh the user's authentication token using a valid refresh token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {refresh_token}"
// @Success 200 {object} map[string]interface{} "auth_token: the new authentication token"
// @Failure 400 {object} map[string]interface{} "error: bad request, token should be provided"
// @Failure 401 {object} map[string]interface{} "error: unauthorized, invalid or expired token"
// @Router /refresh-token [post]
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
