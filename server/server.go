package server

import (
	"file-server/conf"
	_ "file-server/docs"
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
	app = fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "",
		AllowMethods: "GET,POST,HEAD,DELETE",
	}))
	app.Use(Timer())
	app.Use(logger.New())
	app.Use(recover.New())

	setupRoutes()
}

func setupRoutes() {
	if app == nil {
		return
	}

	app.Get("/swagger/*", swagger.HandlerDefault) // default
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		c.JSON(fiber.Map{
			"message": "🐣",
		})
		return c.Next()
	})
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
