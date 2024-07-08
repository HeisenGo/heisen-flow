package http

import (
	"fmt"
	"log"
	"server/config"
	"server/service"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.ServerConfig, app *service.AppContainer) {
	fiberApp := fiber.New()

	// register global routes
	fiberApp.Use(swagger.New())
	// registering users APIs

	// run server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HTTPPort)))
}
