package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/swaggo/fiber-swagger"
	"log"
	"os"
	"path/filepath"
	"server/api/http/handlers"
	"server/api/http/middlewares"
	"server/config"
	_ "server/docs"
	"server/pkg/adapters"
	"server/service"
)

func Run(cfg config.Server, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	createGroupLogger := loggerSetup(fiberApp)

	// register global routes
	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)
	registerGlobalRoutes(api, app, createGroupLogger("global"))
	secret := []byte(cfg.TokenSecret)
	registerBoardRoutes(api, app, secret, createGroupLogger("boards"))
	registerTaskRoutes(api, app, secret, createGroupLogger("tasks"))
	registerColumnRoutes(api, app, secret, createGroupLogger("columns"))
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HTTPPort)))
}

func registerGlobalRoutes(router fiber.Router, app *service.AppContainer, loggerMiddleWare fiber.Handler) {
	router.Use(loggerMiddleWare)
	router.Post("/register", handlers.RegisterUser(app.AuthService()))
	router.Post("/login", handlers.LoginUser(app.AuthService()))
	router.Get("/refresh", handlers.RefreshToken(app.AuthService()))
}

func userRoleChecker() fiber.Handler {
	return middlewares.RoleChecker("user")
}

func registerBoardRoutes(router fiber.Router, app *service.AppContainer, secret []byte, loggerMiddleWare fiber.Handler) {
	router = router.Group("/boards")
	router.Use(loggerMiddleWare)

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

func registerTaskRoutes(router fiber.Router, app *service.AppContainer, secret []byte, loggerMiddleWare fiber.Handler) {
	router = router.Group("/tasks")
	router.Use(loggerMiddleWare)

	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.CreateTask(app.TaskServiceFromCtx),
	)

	router.Get("/:taskID",
		middlewares.Auth(secret),
		handlers.GetFullTaskByID(app.TaskService()),
	)

	router.Patch("/:taskID",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.UpdateTaskColumnByID(app.TaskServiceFromCtx),
	)

	router.Post("/dependency",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.AddDependency(app.TaskServiceFromCtx),
	)
}

func registerColumnRoutes(router fiber.Router, app *service.AppContainer, secret []byte, loggerMiddleWare fiber.Handler) {
	router = router.Group("/columns")
	router.Use(loggerMiddleWare)
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

func loggerSetup(app *fiber.App) func(groupName string) fiber.Handler {

	// Create the logs directory if it does not exist
	logDir := "./logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("error creating logs directory: %v", err)
	}

	// Common format for console logging
	consoleLoggerConfig := logger.Config{
		Format:     "${time} [${ip}]:${port} ${status} - ${method} ${path} - ${latency} ${bytesSent} ${bytesReceived} ${userAgent}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}
	app.Use(logger.New(consoleLoggerConfig))

	// Function to create a logger middleware with separate log file
	createGroupLogger := func(groupName string) fiber.Handler {
		logFilePath := filepath.Join(logDir, groupName+".log")
		file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		return logger.New(logger.Config{
			Format:     "${time} [${ip}]:${port} ${status} - ${method} ${path} - ${latency} ${bytesSent} ${bytesReceived} ${userAgent}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			TimeZone:   "Local",
			Output:     file,
		})
	}
	return createGroupLogger
}
