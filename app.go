package main

import (
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	//database.ConnectMongoDB()
	app := fiber.New()

	app.Use(
		cors.New(cors.Config{
			AllowCredentials: true,
		}),
	)

	routes.AllRoutes(app)
	//routes.ProtectedRoutes(app)

	app.Listen(":8000")
}
