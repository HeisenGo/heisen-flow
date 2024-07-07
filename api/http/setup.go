package http

import (
	"fmt"
	"log"
	"server/api/http/handlers"
	"server/api/http/middlewares"
	"server/config"
	"server/service"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Server, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	// register global routes
	registerGlobalRoutes(api, app)

	//secret := []byte(cfg.TokenSecret)

	// registering users APIs
	//registerUsersAPI(api, app.UserService(), secret)

	// run server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HTTPPort)))
}

//func registerUsersAPI(router fiber.Router, _ *service.UserService, secret []byte) {
//	userGroup := router.Group("/users", middlewares.Auth(secret), middlewares.RoleChecker("user", "admin"))
//
//	userGroup.Get("/claims", func(c *fiber.Ctx) error {
//		claims := c.Locals(jwt.UserClaimKey).(*jwt.UserClaims)
//
//		return c.JSON(map[string]any{
//			"user_id": claims.UserID,
//			"role":    claims.Role,
//		})
//	})
//}

func registerGlobalRoutes(router fiber.Router, app *service.AppContainer) {
	router.Post("/register", handlers.RegisterUser(app.AuthService()))
	router.Post("/login", handlers.LoginUser(app.AuthService()))
	router.Get("/refresh", handlers.RefreshToken(app.AuthService()))
}

func userRoleChecker() fiber.Handler {
	return middlewares.RoleChecker("user")
}
