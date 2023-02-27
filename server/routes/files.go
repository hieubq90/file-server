package routes

import (
	"file-server/server/handlers"

	"github.com/gofiber/fiber/v2"
)

// FilesRouter is the Router for GoFiber App
func FilesRouter(app fiber.Router) {
	app.Post("/:project/:folder/files", handlers.UploadFile())
}
