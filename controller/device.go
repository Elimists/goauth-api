package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetDevices(c *fiber.Ctx) error {

	//RETURN_LIMIT := 12
	params := c.AllParams()
	pagenum, _ := strconv.Atoi(params["pagenumber"])
	if pagenum == 0 {
		pagenum = 1
	}

	//offset := (pagenum - 1) * RETURN_LIMIT

	var devices []*models.Device

	//database.DB.Offset(offset).Limit(RETURN_LIMIT).Model(&models.Device{Stage: "public"}).Preload("Images").Find(&devices)

	database.DB.Select("id", "name", "author", "time_to_complete", "material_cost", "created_at", "updated_at").Where("stage = ?", "public").Find(&devices).Preload("Images")
	return c.JSON(&devices)
}

func GetDevice(c *fiber.Ctx) error {

	params := c.AllParams()
	id := params["id"]
	var device models.Device

	if err := database.DB.Where("id = ?", &id).Preload("Capabilities").Preload("Disabilities").Preload("Usages").First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Device not found!"}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	return c.JSON(&device)
}

func GetDeviceComments(c *fiber.Ctx) error {

	RETURN_LIMIT := 8
	params := c.AllParams()
	id := params["id"]
	pagenum, _ := strconv.Atoi(params["pagenumber"])
	if pagenum == 0 {
		pagenum = 1
	}

	offset := (pagenum - 1) * RETURN_LIMIT

	var comments []*models.Review

	if err := database.DB.Offset(offset).Limit(RETURN_LIMIT).Where("id = ?", &id).Preload("Capabilities").Preload("Disabilities").Preload("Usages").Find(&comments).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "No comments found for the device"}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	return c.JSON(&comments)
}

func AddDevice(c *fiber.Ctx) error {
	/*
		cookie := c.Cookies("mmc_cookie")

		authpacket := Checkauth(cookie)

		if authpacket.Error {
			return c.Status(fiber.StatusUnauthorized).JSON(authpacket)
		}

		var device models.Device
		if err := c.BodyParser(&device); err != nil {
			rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
			return c.Status(fiber.StatusNotAcceptable).JSON(rp)
		}

		device.UUID = uuid.NewString()
		device.UrlSafeName = strings.ToLower(strings.ReplaceAll(device.Name, " ", "_"))
		device.Author = authpacket.Claim.Email

		if err := database.DB.Create(&device).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				rp := models.ResponsePacket{Error: true, Code: "create_error", Message: "Device with that name already exists. Please use another name"}
				return c.Status(fiber.StatusNotAcceptable).JSON(rp)
			}
			rp := models.ResponsePacket{Error: true, Code: "create_error", Message: "Unable to create device. Internal error"}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}
	*/
	rp := models.ResponsePacket{Error: false, Code: "create_success", Message: "Device created succesfully"}
	return c.Status(fiber.StatusAccepted).JSON(rp)
}

func SaveFile(c *fiber.Ctx) error {

	file, err := c.FormFile("relatedfile")
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "no_attachment", Message: "No pictures attached" + err.Error()}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	filext := strings.Split(file.Filename, ".")[1]
	if filext != "rar" && filext != "zip" {
		rp := models.ResponsePacket{Error: true, Code: "unsupported_type", Message: "Unsupported file type. Make sure all files are packaged as zip or rar."}
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(rp)
	}

	fmt.Println(file.Size)
	return c.JSON("")
}

func SaveImages(c *fiber.Ctx) error {

	form, err := c.MultipartForm()
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "no_attachment", Message: "No images attached. Not acceptable"}
		return c.JSON(rp)
	}

	files := form.File["images"]
	if len(files) > 5 {
		rp := models.ResponsePacket{Error: true, Code: "quantity_exceeded", Message: "You can only upload 5 images"}
		return c.JSON(rp)
	}

	return c.JSON("")
}
