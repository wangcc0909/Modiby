package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.peaut.limit/filearound/router"
)

func Run() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Asura",
	})
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	router.Route(app)
	app.Listen(":3000")
}
