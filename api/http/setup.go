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
	registerColumnRoutes(api, app, secret)
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HTTPPort)))
}

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

	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.CreateUserBoard(app.BoardServiceFromCtx),
	)
	router.Get("/my-boards",
		middlewares.Auth(secret),
		handlers.GetUserBoards(app.BoardService()),
	)
	router.Get("/publics",
		middlewares.Auth(secret),
		handlers.GetPublicBoards(app.BoardService()),
	)
	router.Get("/:boardID",
		middlewares.Auth(secret),
		handlers.GetFullBoardByID(app.BoardService()),
	)

	router.Delete("/:boardID",
		middlewares.Auth(secret),
		handlers.DeleteBoard(app.BoardService()),
	)

	router.Post("/invite", middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.InviteToBoard(app.BoardServiceFromCtx))
}

func registerTaskRoutes(router fiber.Router, app *service.AppContainer, secret []byte) {
	router = router.Group("/tasks")

	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.CreateTask(app.TaskServiceFromCtx),
	)

	router.Get("/:taskID",
		middlewares.Auth(secret),
		handlers.GetFullTaskByID(app.TaskService()),
	)

	router.Post("/dependency",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.AddDependency(app.TaskServiceFromCtx),
	)
}

func registerColumnRoutes(router fiber.Router, app *service.AppContainer, secret []byte) {
	router = router.Group("/columns")
	router.Post("",
		middlewares.Auth(secret),
		userRoleChecker(),
		handlers.CreateColumns(app.ColumnService()),
	)
	router.Delete("/:columnID",
		middlewares.Auth(secret),
		handlers.DeleteColumn(app.ColumnService()),
	)

	router.Put("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.ReorderColumns(app.ColumnServiceFromCtx),
	)
}
