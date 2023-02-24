package server

import (
	"file-server/conf"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var app *fiber.App

func InitServer() {
	app = fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "",
		AllowMethods: "GET,POST,HEAD,DELETE",
	}))
	app.Use(logger.New())

	setupRoutes()
}

func setupRoutes() {

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

}
