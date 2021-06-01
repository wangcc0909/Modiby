package router

import (
	"github.com/gofiber/fiber/v2"
	"github.peaut.limit/filearound/service"
)

func Route(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	{
		v1.Get("/file/list", service.GetFileList)
		v1.Post("/file/upload", service.UploadFile)
		v1.Get("/download/:filename", service.Download)
	}
}
