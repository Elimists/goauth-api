package routes

import (
	"fmt"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/middleware"
	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app *fiber.App) {

	version := "/api/v2"
	/*AUTH Routes*/
	app.Post(fmt.Sprintf("%s/verify", version), controller.VerifyEmail)
	app.Post(fmt.Sprintf("%s/register", version), middleware.Limiter(5, 40), controller.Register)
	app.Post(fmt.Sprintf("%s/login", version), middleware.Limiter(6, 45), controller.Login)
	app.Post(fmt.Sprintf("%s/resetpassword", version), middleware.Limiter(6, 45), controller.ResetPassword)

	/*USER Routes*/
	app.Get(fmt.Sprintf("%s/getuser", version), controller.GetUser)
	//app.Get(fmt.Sprintf("api/%s/getprofilepic", version), controller.GetProfilePic)

	/*DEVICE Routes*/
	app.Get(fmt.Sprintf("%s/getdevices/:pagenumber", version), controller.GetDevices)
	app.Get(fmt.Sprintf("%s/getdevice/:id", version), controller.GetDevice)
	//temporary
	//app.Post("api/v2/savefile", controller.SaveFile)

	app.Get(fmt.Sprintf("%s/getallusers", version), controller.GetAllUsers)

	// PROTECTED ROUTES
	//app.Get(fmt.Sprintf("%s/getallusers", version), middleware.Protected(), controller.GetAllUsers)
	app.Patch(fmt.Sprintf("%s/updateuser", version), middleware.Protected(), middleware.Limiter(6, 60), controller.UpdateUser)
	app.Post(fmt.Sprintf("%s/createdevice", version), middleware.Protected(), controller.AddDevice)
	//app.Patch(fmt.Sprintf("%s/uploadpic", version), middleware.Protected(), controller.UpdateProfilePic)
	app.Post(fmt.Sprintf("%s/updatepassword", version), middleware.Protected(), middleware.Limiter(6, 45), controller.UpdatePassword)

}
