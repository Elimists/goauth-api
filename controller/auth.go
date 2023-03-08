package controller

import (
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if !emailIsValid(data["email"]) {
		rp := models.ResponsePacket{Error: true, Code: "invalid_email", Message: "Email is not a valid type."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if !passwordIsValid(data["password"]) {
		rp := models.ResponsePacket{Error: true, Code: "invalid_password", Message: "Password is not strong enough."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)
	datetime := time.Now().Unix()

	auth := models.Auth{
		Email:        data["email"],
		Password:     password,
		Verified:     false,
		Privilege:    9, // General user.
		RegisteredOn: uint(datetime),
		LastLoggedIn: uint(datetime),
	}

	authErr := database.DB.Create(&auth).Error

	if authErr != nil {
		if strings.Contains(authErr.Error(), "Duplicate entry") {
			rp := models.ResponsePacket{Error: true, Code: "duplicate_email", Message: "Email already exists!"}
			return c.Status(fiber.StatusNotAcceptable).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error. Could not register user."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	verification := models.Verification{
		Code:    generateVerificationCode(),
		Email:   data["email"],
		Expires: uint(time.Now().Add(time.Minute * 30).Unix()),
	}

	verifyErr := database.DB.Create(&verification).Error
	if verifyErr != nil {
		database.DB.Delete(&auth).Where("email = ?", data["email"]) // delete the row created prior in Auth table
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error. Could not register user."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	//Send verification email here
	emailSendingError := SendVerificationCode(data["email"], verification.Code)

	if emailSendingError {
		database.DB.Where("email = ?", data["email"]).Delete(&verification)
		database.DB.Where("email = ?", data["email"]).Delete(&auth)
		rp := models.ResponsePacket{Error: false, Code: "email_sending_error", Message: "Unable to send email.DB Row rollbacked."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	rp := models.ResponsePacket{Error: false, Code: "user_registered", Message: "User successfully registered."}
	return c.Status(fiber.StatusCreated).JSON(rp)
}

// Login route method
//
// On success, sets X-Maker-Token:{token} in header.
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var auth models.Auth

	if err := database.DB.Where("email = ?", data["email"]).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: fiber.StatusNotFound, Message: "Account not found!"}
		}
	}

	if !auth.Verified {
		rp := models.ResponsePacket{Error: true, Code: "email_unverified", Message: "User is not verfied."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	expiry := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	if data["longerlogin"] == "yes" {
		expiry = jwt.NewNumericDate(time.Now().Add(240 * time.Hour))
	}

	if err := bcrypt.CompareHashAndPassword(auth.Password, []byte(data["password"])); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "incorrect_password", Message: "Password is not correct"}
		return c.Status(fiber.StatusBadRequest).JSON(rp)
	}

	claims := jwt.MapClaims{
		"email":     auth.Email,
		"verified":  auth.Verified,
		"privilege": auth.Privilege,
		"exp":       expiry,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not sign token."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	database.DB.Model(&auth).Where("email = ?", data["email"]).Update("last_logged_in", time.Now().Unix()) // update the last logged in datetime

	c.Append("X-Maker-Token", signedToken)

	rp := models.ResponsePacket{Error: false, Code: "successfull", Message: "Login successfull"}
	return c.Status(fiber.StatusOK).JSON(rp)
}

// Email Verification Route
//
// TODO Needs to me changed to link verification method
func VerifyEmail(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var verification models.Verification

	if err := database.DB.Where("email = ?", data["email"]).First(&verification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Verification code for this user does not exsist"}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
	}

	if uint(time.Now().Unix()) > verification.Expires {
		rp := models.ResponsePacket{Error: true, Code: "expired", Message: "Verfication time frame has expired."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if verification.Code != data["code"] {
		rp := models.ResponsePacket{Error: true, Code: "code_mismatch", Message: "Verification code does not match."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if err := database.DB.Model(&verification).Where("email = ?", data["email"]).Updates(map[string]interface{}{"verified": true}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Could not set verification status to True"}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}
	}

	rp := models.ResponsePacket{Error: false, Code: "verified", Message: "Verification of email successfull!"}
	return c.Status(fiber.StatusAccepted).JSON(rp)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if len(data["oldpassword"]) <= 0 || len(data["newpassword"]) <= 0 {
		rp := models.ResponsePacket{Error: true, Code: "missing_data", Message: "Form is missing required data!"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	//TODO - Implement changing password functionality
	if !passwordIsValid(data["oldpassword"]) || !passwordIsValid(data["newpassword"]) {
		rp := models.ResponsePacket{Error: true, Code: "invalid_password", Message: "Password is not strong enough."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthroized"})

}

func ResetPassword(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthroized"})
}

//__________________________________________________________________________
/*
* HELPER FUNCTIONS
 */
func emailIsValid(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(s)
}

func passwordIsValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 10101
	max := 99999
	return strconv.Itoa((rand.Intn(max-min+1) + min))
}

// Send the verification code to the user.
func SendVerificationCode(email string, verificationCode string) bool {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "231c63d58c7571", "15065dc065bf4c", "sandbox.smtp.mailtrap.io")

	to := []string{email}
	subject := "Subject: Makers Verification Code\n"
	from := "maker@example.com"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<div  style="font-size:20px; font-family: Arial, serif;">
				<p>Here is your verification code: <code style="font-weight: bold;">%s</code></p>
			</div>
		</html>
		`, verificationCode)
	msg := []byte(subject + mime + body)

	err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", auth, from, to, msg)
	if err != nil {
		return true
	}

	return false
}
