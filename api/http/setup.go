package http

import (
	"fmt"
	"log"
	"server/api/http/handlers"
	"server/api/http/middlewares"
	"server/config"
	"server/pkg/adapters"
	"server/service"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Server, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	// register global routes
	registerGlobalRoutes(api, app)
	secret := []byte(cfg.TokenSecret)
	registerBoardRoutes(api, app, secret)
	registerTaskRoutes(api, app, secret)

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

func registerBoardRoutes(router fiber.Router, app *service.AppContainer, secret []byte) {
	router = router.Group("/boards")

	//router.Get("", middlerwares.Auth(secret), userRoleChecker(), handlers.UserOrders(app.OrderService()))

	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		//userRoleChecker(),
		handlers.CreateUserBoard(app.BoardServiceFromCtx),
	)

	router.Post("/invite", middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		//	userRoleChecker(),
		handlers.InviteToBoard(app.BoardServiceFromCtx))
}

func registerTaskRoutes(router fiber.Router, app *service.AppContainer, secret []byte) {
	router = router.Group("/tasks")

	//router.Get("", middlerwares.Auth(secret), userRoleChecker(), handlers.UserOrders(app.OrderService()))

	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		//userRoleChecker(),
		handlers.CreateTask(app.TaskServiceFromCtx),
	)
}
