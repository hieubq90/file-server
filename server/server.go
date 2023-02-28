package server

import (
	"encoding/json"
	"file-server/conf"
	_ "file-server/docs"
	"file-server/mio"
	"file-server/server/routes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"time"
)

var app *fiber.App

// Tính toán thời gian phản hồi
func Timer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// start timer
		start := time.Now()
		// next routes
		err := c.Next()
		// stop timer
		stop := time.Now()
		// Do something with response
		c.Append("Server-Timing", fmt.Sprintf("app;dur=%v", stop.Sub(start).String()))
		// return stack error if exist
		return err
	}
}

func InitServer() {
	serverConfig := conf.GetServerConfig()
	app = fiber.New(fiber.Config{
		AppName:                      serverConfig["name"],
		UnescapePath:                 true,
		BodyLimit:                    50 * 1024 * 1024,
		DisablePreParseMultipartForm: false,
		JSONEncoder:                  json.Marshal,
		JSONDecoder:                  json.Unmarshal,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "",
		AllowMethods: "GET,POST,HEAD,DELETE",
	}))
	app.Use(Timer())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Hooks().OnListen(func() error {
		minioConfig := conf.GetMinioConfig()
		err := mio.CreateBucket(minioConfig["bucket"])
		return err
	})

	setupRoutes()
}

func setupRoutes() {
	if app == nil {
		return
	}

	app.Get("/swagger/*", swagger.HandlerDefault) // default
	api := app.Group("/api")

	routes.FilesRouter(api)
}

func StartServer() (err error) {
	//appConfig := conf.GetAppConfig()
	//endpoint := conf.GetApplicationName()
	serverConfig := conf.GetServerConfig()
	fmt.Printf("Khởi chạy ứng dụng %s", serverConfig["name"])
	err = app.Listen(serverConfig["endpoint"])
	return err
}

func Shutdown() {
	app.Shutdown()
}
