package middlewares

import (
	"errors"
	"server/api/http/handlers"
	"server/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.GetReqHeaders()["Authorization"][0]

		if len(h) == 0 {
			return handlers.SendError(c, errors.New("authorization token not specified"), fiber.StatusUnauthorized)
		}

		pureToken := strings.Split(h, " ")[1]
		claims, err := jwt.ParseToken(pureToken, secret)
		if err != nil {
			return handlers.SendError(c, err, fiber.StatusUnauthorized)
		}

		c.Locals(jwt.UserClaimKey, claims)

		return c.Next()
	}
}

func RoleChecker(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals(jwt.UserClaimKey).(*jwt.UserClaims)
		hasAccess := false
		for _, role := range roles {
			if claims.Role == role {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return handlers.SendError(c, errors.New("you don't have access to this section"), fiber.StatusForbidden)
		}

		return c.Next()
	}
}
