package controller

import (
	"errors"
	"math/rand"
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

var SECRET_KEY = []byte("K7yx09lpbR")

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

	user := models.User{
		Email: data["email"],
	}
	if userErr := database.DB.Create(&user).Error; userErr != nil {
		database.DB.Delete(&auth).Where("email = ?", data["email"])
		database.DB.Delete(&verification).Where("email = ?", data["email"])
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error. Could not register user."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	rp := models.ResponsePacket{Error: false, Code: "user_registered", Message: "User successfully registered."}
	return c.Status(fiber.StatusCreated).JSON(rp)
}

/*
* LOGIN
* Handles the login flow.
* @returns:
 */
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var auth models.Auth

	if err := database.DB.Where("email = ?", data["email"]).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Email not found!"}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
	}

	expiry := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	if data["longerlogin"] == "yes" {
		expiry = jwt.NewNumericDate(time.Now().Add(240 * time.Hour))
	}

	if err := bcrypt.CompareHashAndPassword(auth.Password, []byte(data["password"])); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "incorrect_password", Message: "Password is not correct"}
		return c.Status(fiber.StatusBadRequest).JSON(rp)
	}

	claims := models.CustomClaims{
		Email:     auth.Email,
		Verified:  auth.Verified,
		Privilege: auth.Privilege,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Mak-Mak-Ch",
			ExpiresAt: expiry,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(SECRET_KEY)

	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not log in." + err.Error()}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	expires := time.Now().Add(time.Hour * 8)
	if data["longexpiry"] == "yes" {
		expires = time.Now().Add(time.Hour * 244) // JWT token expires after 10 days.
	}

	cookie := fiber.Cookie{
		Name:     "mmc_cookie",
		Value:    ss,
		Expires:  expires,
		HTTPOnly: true,
	}

	database.DB.Model(&auth).Where("email = ?", data["email"]).Update("last_logged_in", time.Now().Unix()) // update the last logged in datetime

	c.Cookie(&cookie)

	rp := models.ResponsePacket{Error: false, Code: "successfull", Message: "Login successfull"}
	return c.Status(fiber.StatusOK).JSON(rp)
}

/*
* VERIFY EMAIL
* Handles the email verification flowError
* @returns: http status and JSON response
 */
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
