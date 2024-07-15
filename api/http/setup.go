package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
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

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New(fiber.Config{
		Views: html.New("./templates", ".html"),
	})
	// Serve static files from the "assets" directory
	fiberApp.Static("/assets", "./assets")

	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	createGroupLogger := loggerSetup(fiberApp)

	registerUIRoutes(fiberApp, app, createGroupLogger("templates"))
	// register global routes
	registerGlobalRoutes(api, app,
		createGroupLogger("global"),
		middlewares.SetupLimiterMiddleware(1, 1, cfg.Redis),
	)
	secret := []byte(cfg.Server.TokenSecret)
	registerBoardRoutes(api, app, secret, createGroupLogger("boards"))
	registerTaskRoutes(api, app, secret, createGroupLogger("tasks"))
	registerColumnRoutes(api, app, secret, createGroupLogger("columns"))
	registerNotificationRoutes(api, app, secret, createGroupLogger("notifs"))
	registerCommentRoutes(api, app, secret, createGroupLogger("comments"))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)))
}

func registerUIRoutes(router fiber.Router, app *service.AppContainer, loggerMiddleWare fiber.Handler) {
	router.Use(loggerMiddleWare)
	router.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "HeisenFlow",
		})
	})
	router.Get("/swagger/*", fiberSwagger.WrapHandler)
	router.Get("/metrics", monitor.New(monitor.Config{Title: "HeisenFlow Metrics Page"}))
}

func registerGlobalRoutes(router fiber.Router, app *service.AppContainer, loggerMiddleWare fiber.Handler, limiterMiddleWare fiber.Handler) {
	router.Use(loggerMiddleWare)
	router.Post("/register", limiterMiddleWare, handlers.RegisterUser(app.AuthService()))
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
		middlewares.SetupCacheMiddleware(5),
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

	router.Patch("/reorder",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.ReorderTasks(app.TaskServiceFromCtx),
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
func registerNotificationRoutes(router fiber.Router, app *service.AppContainer, secret []byte, loggerMiddleWare fiber.Handler) {
	router = router.Group("/notifications")
	router.Use(loggerMiddleWare)
	router.Get("", middlewares.Auth(secret), handlers.GetNotifications(app.NotificationService()))
	router.Patch("/read/:notifID", middlewares.Auth(secret), handlers.UpdateNotifications(app.NotificationService()))
}

func registerCommentRoutes(router fiber.Router, app *service.AppContainer, secret []byte, loggerMiddleWare fiber.Handler) {
	router = router.Group("/comments")
	router.Use(loggerMiddleWare)

	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		middlewares.Auth(secret),
		handlers.CreateUserComment(app.CommentServiceFromCtx),
	)
}
