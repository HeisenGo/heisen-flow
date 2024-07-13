package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"time"
)

// SetupCacheMiddleware Define a function to configure cache middleware
func SetupCacheMiddleware(expMinutes int) fiber.Handler {
	exp := time.Duration(expMinutes)
	return cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   exp * time.Minute,
		CacheControl: true,
	})
}
