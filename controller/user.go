package controller

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GetUser(c *fiber.Ctx) error {

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	var user models.User

	if err := database.DB.Where("email = ?", claims["email"].(string)).First(&user).Error; err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error"}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	return c.Status(fiber.StatusOK).JSON(&user)
}

func GetAllUsers(c *fiber.Ctx) error {

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if claims["privilege"].(float64) > 3 {
		rp := models.ResponsePacket{Error: true, Code: "insufficient_privileges", Message: "You do not have sufficient privilege to perform this action!"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	RETURN_LIMIT := 30
	params := c.AllParams()
	pagenum, _ := strconv.Atoi(params["pagenumber"])
	if pagenum == 0 {
		pagenum = 1
	}

	offset := (pagenum - 1) * RETURN_LIMIT

	var users []*models.User

	database.DB.Offset(offset).Limit(RETURN_LIMIT).Find(&users)

	return c.Status(fiber.StatusAccepted).JSON(&users)

}

func UpdateUser(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	user := models.User{
		Name:         data["name"],
		Education:    data["education"],
		Occupation:   data["occupation"],
		Organization: data["organization"],
		Bio:          data["bio"],
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if err := database.DB.Where("email = ?", claims["email"].(string)).Updates(&user).Error; err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error. Could not update user info"}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	rp := models.ResponsePacket{Error: false, Code: "update_successfull", Message: "User info successfully updated."}
	return c.Status(fiber.StatusCreated).JSON(rp)
}

func UpdateProfilePic(c *fiber.Ctx) error {

	file, err := c.FormFile("profilepic")
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "no_attachment", Message: "No pictures attached" + err.Error()}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	//fmt.Println(file.Size)
	if file.Size > 75000 {
		rp := models.ResponsePacket{Error: true, Code: "too_large", Message: "Profile picture is too large. Image must be less than 75kb"}
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(rp)
	}

	id := uuid.New()
	filename := strings.Replace(id.String(), "-", "", -1)
	filext := strings.Split(file.Filename, ".")[1]
	if filext != "jpg" && filext != "png" {
		rp := models.ResponsePacket{Error: true, Code: "unsupported_type", Message: "Unsupported image type."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var user models.User

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if err := database.DB.Where("email = ?", claims["email"].(string)).First(&user).Error; err != nil {
		rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "User not found."}
		return c.Status(fiber.StatusNotFound).JSON(rp)
	}

	if len(user.Picture) != 0 {
		if err := os.Remove(fmt.Sprintf("./uploads/users/profilepics/%s", user.Picture)); err != nil {
			database.DB.Model(&user).Where("email = ?", claims["email"].(string)).Update("picture", "")
			rp := models.ResponsePacket{Error: true, Code: "unable_to_delete", Message: "Unable to update profile pic."}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}
		if err := database.DB.Model(&user).Where("email = ?", claims["email"].(string)).Update("picture", "").Error; err != nil {
			rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not delete existing image string from database" + err.Error()}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}
	}

	image := fmt.Sprintf("%s.%s", filename, filext)

	if err := database.DB.Model(&user).Where("email = ?", claims["email"].(string)).Update("picture", image).Error; err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not save image string to database" + err.Error()}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	if err := c.SaveFile(file, fmt.Sprintf("./uploads/users/profilepics/%s", image)); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal error. Unable to save profile pic. Line 156 - " + err.Error()}
		database.DB.Model(&user).Where("email = ?", claims["email"].(string)).Update("picture", "")
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	rp := models.ResponsePacket{Error: false, Code: "save_successfull", Message: "Profile image updated successfully."}
	return c.Status(fiber.StatusAccepted).JSON(rp)
}

func GetProfilePic(c *fiber.Ctx) error {

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	var user models.User

	if err := database.DB.Where("email = ?", claims["email"].(string)).First(&user).Error; err != nil {
		rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "User not found."}
		return c.Status(fiber.StatusNotFound).JSON(rp)
	}

	return c.SendFile(fmt.Sprintf("./uploads/users/profilepics/%s", user.Picture), true)
}

func HandleTokenCheck(c *fiber.Ctx) (*jwt.Token, error) {
	cookie := c.Cookies("mmc_cookie")
	token, err := jwt.ParseWithClaims(cookie, &models.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		})

	return token, err
}
