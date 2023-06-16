package controller

import (
	"strconv"

	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

	RETURN_LIMIT := 30
	params := c.AllParams()
	pagenum, _ := strconv.Atoi(params["pagenumber"])
	if pagenum == 0 {
		pagenum = 1
	}

	offset := (pagenum - 1) * RETURN_LIMIT

	var users []*models.User

	database.DB.Offset(offset).Limit(RETURN_LIMIT).Preload("UserDetails").Find(&users)

	return c.Status(fiber.StatusAccepted).JSON(&users)

}

func UpdateUser(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	user := models.UserDetails{
		FirstName:    data["firstname"],
		LastName:     data["lastname"],
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

func GetAddress(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&models.ResponsePacket{Error: false, Code: "success", Message: "Address added successfully"})
}

func AddAddress(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&models.ResponsePacket{Error: false, Code: "success", Message: "Address added successfully"})
}

func UpdateAddress(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&models.ResponsePacket{Error: false, Code: "success", Message: "Address added successfully"})
}

func DeleteAddress(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&models.ResponsePacket{Error: false, Code: "success", Message: "Address added successfully"})
}
