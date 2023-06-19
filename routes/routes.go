package routes

import (
	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/middleware"
	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/html/home.html")
	})
	/*AUTH Routes*/
	app.Get("/verify/:email/:verificationCode", controller.VerifyEmail)
	app.Get("/register", controller.ShowRegistrationForm)
	app.Post("/register", middleware.Limiter(14, 60), controller.Register)
	app.Post("/login", middleware.Limiter(6, 45), controller.Login)
	app.Post("/logout", middleware.Limiter(6, 45), controller.Logout)
	app.Post("/resetpassword", middleware.Limiter(6, 45), controller.ResetPassword)

	/*USER Routes*/
	app.Get("/getuser", controller.GetUser)
	//app.Get("api//getprofilepic", controller.GetProfilePic)

	/*DEVICE Routes*/
	app.Get("/getdevices/:pagenumber", controller.GetDevices)
	app.Get("/getdevice/:id", controller.GetDevice)
	//temporary
	//app.Post("api/v2/savefile", controller.SaveFile)

	//app.Get("/getallusers", controller.GetAllUsers)

	// PROTECTED ROUTES
	app.Get("/users", middleware.Protected(), controller.GetAllUsers)
	app.Get("/users/:id", middleware.Protected(), middleware.Limiter(6, 60), controller.GetUser)
	app.Patch("/users/:id", middleware.Protected(), middleware.Limiter(6, 60), controller.UpdateUser)

	app.Get("/users/:id/address", middleware.Protected(), controller.GetAddress)
	app.Post("/users/:id/address", middleware.Protected(), controller.AddAddress)
	app.Patch("/users/:id/address/:id", middleware.Protected(), controller.UpdateAddress)
	app.Delete("/users/:id/address/:id", middleware.Protected(), controller.DeleteAddress)

	//app.Patch("/uploadpic", middleware.Protected(), controller.UpdateProfilePic)
	app.Post("/updatepassword", middleware.Protected(), middleware.Limiter(6, 45), controller.UpdatePassword)

}
