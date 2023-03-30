package controller

import (
	"html/template"
	"strconv"

	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Name  string
	Email string
}

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

	/*
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

		database.DB.Offset(offset).Limit(RETURN_LIMIT).Preload("User").Find(&users)

		return c.Status(fiber.StatusAccepted).JSON(&users)
	*/

	users := []User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
		{Name: "David", Email: "david@example.com"},
		{Name: "Emily", Email: "emily@example.com"},
		{Name: "Frank", Email: "frank@example.com"},
		{Name: "Grace", Email: "grace@example.com"},
		{Name: "Henry", Email: "henry@example.com"},
		{Name: "Ivy", Email: "ivy@example.com"},
		{Name: "Jack", Email: "jack@example.com"},
	}

	// Parse the template
	t, err := template.ParseFiles("./static/html/users.html")
	if err != nil {
		return err
	}

	// Determine the page number and page size
	page := 1
	pageSize := 3
	if pageParam := c.Query("page"); pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(users) {
		endIndex = len(users)
	}

	// Filter the user data to the current page
	filteredUsers := users[startIndex:endIndex]

	// Render the template with the user data and pagination links
	err = t.Execute(c.Response().BodyWriter(), struct {
		Users       []User
		CurrentPage int
		TotalPages  int
		PrevPage    int
		NextPage    int
	}{
		Users:       filteredUsers,
		CurrentPage: page,
		TotalPages:  (len(users) + pageSize - 1) / pageSize,
		PrevPage:    page - 1,
		NextPage:    page + 1,
	})
	if err != nil {
		return err
	}

	// Set the Content-Type header to text/html
	c.Set("Content-Type", "text/html")

	return nil
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
