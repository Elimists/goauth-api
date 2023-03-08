package routes

import (
	"fmt"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/middleware"
	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app *fiber.App) {

	version := "v2"
	/*AUTH Routes*/
	app.Post(fmt.Sprintf("/api/%s/verify", version), controller.VerifyEmail)
	app.Post(fmt.Sprintf("/api/%s/register", version), controller.Register)
	app.Post(fmt.Sprintf("/api/%s/login", version), middleware.Limiter(6, 45), controller.Login)
	app.Post(fmt.Sprintf("/api/%s/resetpassword", version), controller.ResetPassword)

	/*USER Routes*/
	app.Get(fmt.Sprintf("api/%s/getuser", version), controller.GetUser)
	app.Get(fmt.Sprintf("api/%s/getprofilepic", version), controller.GetProfilePic)

	/*DEVICE Routes*/
	app.Get(fmt.Sprintf("api/%s/getdevices/:pagenumber", version), controller.GetDevices)
	app.Get(fmt.Sprintf("api/%s/getdevice/:id", version), controller.GetDevice)
	//temporary
	//app.Post("api/v2/savefile", controller.SaveFile)

	// PROTECTED ROUTES
	app.Get(fmt.Sprintf("api/%s/getallusers", version), middleware.Protected(), controller.GetAllUsers)
	app.Patch(fmt.Sprintf("/api/%s/updateuser", version), middleware.Protected(), middleware.Limiter(6, 60), controller.UpdateUser)
	app.Post(fmt.Sprintf("api/%s/createdevice", version), middleware.Protected(), controller.AddDevice)
	app.Patch(fmt.Sprintf("api/%s/uploadpic", version), middleware.Protected(), controller.UpdateProfilePic)

}
