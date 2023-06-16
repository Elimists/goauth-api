package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

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
		KeyLookup:      fmt.Sprintf("header:X-%s-CSRF-Token", os.Getenv("API_NAME")),
		CookieName:     fmt.Sprintf("%s_csrf", os.Getenv("API_NAME")),
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
	}))

	routes.AllRoutes(app)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("API_PORT")))

}
