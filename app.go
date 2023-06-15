package main

import (
	"time"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

func main() {
	database.Connect()
	go controller.EmailVerificationWorker()
	app := fiber.New()

	app.Use(
		cors.New(cors.Config{
			AllowCredentials: true,
		}),
	)

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Maker-Csrf-Token",
		CookieName:     "csrf_1",
		CookieSameSite: "Lax",
		Expiration:     15 * time.Minute,
	}))

	routes.AllRoutes(app)
	//routes.ProtectedRoutes(app)

	app.Listen(":8000")
}
