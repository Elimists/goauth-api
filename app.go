package main

import (
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	database.Connect()
	app := fiber.New()
	app.Use(
		cors.New(cors.Config{
			AllowCredentials: true,
		}),
		limiter.New())
	routes.Setup(app)
	//http.HandleFunc("/api/v2/upload", controller.UploadImage)
	app.Listen(":8000")
}
