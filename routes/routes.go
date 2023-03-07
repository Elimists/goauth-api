package routes

import (
	"github.com/Elimists/go-app/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	/*AUTH Routes*/
	app.Post("/api/v2/verify", controller.VerifyEmail)
	app.Post("/api/v2/register", controller.Register)
	app.Post("/api/v2/login", controller.Login)

	/*USER Routes*/
	app.Patch("/api/v2/updateuser", controller.UpdateUser)
	app.Get("api/v2/getuser", controller.GetUser)
	app.Get("api/v2/getallusers", controller.GetAllUsers)
	app.Patch("api/v2/uploadpic", controller.UpdateProfilePic)
	app.Get("api/v2/getprofilepic", controller.GetProfilePic)

	/*DEVICE Routes*/
	app.Get("api/v2/getdevices/:pagenumber", controller.GetDevices)
	app.Get("api/v2/getdevice/:id", controller.GetDevice)
	app.Post("api/v2/createdevice", controller.AddDevice)
	//temporary
	//app.Post("api/v2/savefile", controller.SaveFile)

}
