package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"time"
)

// SetupLimiterMiddleware Define a function to configure limiter middleware
func SetupLimiterMiddleware(durationMinutes int, max int) fiber.Handler {
	// Create limiter middleware with customized settings
	exp := time.Duration(durationMinutes)
	return limiter.New(limiter.Config{
		// Function to determine if request should be limited (optional)
		//Next: func(c *fiber.Ctx) bool {
		//	// Example: Limit requests only for localhost (disable for local testing)
		//	return c.IP() == "127.0.0.1"
		//},
		Max:        max,
		Expiration: exp * time.Minute,
		Storage:    redis.New(),
	})
}
