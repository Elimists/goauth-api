package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	"github.com/eapache/channels"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var emailQueue = channels.NewInfiniteChannel()

// Register a new user.
//
// @Summary Register a new user.
// @Description Register a new user and send them a verification email.
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
	verificationCode := generateVerificationCode()
	auth := models.User{
		Email:              data["email"],
		Password:           password,
		Privilege:          9, // General user.
		Verified:           false,
		VerificationCode:   verificationCode,
		VerificationExpiry: uint(time.Now().Add(time.Minute * 30).Unix()),
		UserDetails: models.UserDetails{
			FirstName: data["firstName"],
			LastName:  data["lastName"],
		},
	}

	userErr := database.DB.Create(&auth).Error

	if userErr != nil {
		if strings.Contains(userErr.Error(), "Duplicate entry") {
			rp := models.ResponsePacket{Error: true, Code: "duplicate_email", Message: "Email already exists!"}
			return c.Status(fiber.StatusNotAcceptable).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error. Could not register user."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(data["email"]))
	encodedVerificationCode := base64.StdEncoding.EncodeToString([]byte(verificationCode))

	verificationLink := fmt.Sprintf("http://localhost:8000/api/v2/verify/%s/%s", encodedEmail, encodedVerificationCode)

	// Add email and verificationCode to inmemory queue.
	body := fmt.Sprintf(`{"email": "%s", "verificationLink": "%s"}`, encodedEmail, verificationLink)
	emailQueue.In() <- body

	rp := models.ResponsePacket{Error: false, Code: "user_registered", Message: "User registered successfully."}
	return c.Status(fiber.StatusCreated).JSON(rp)
}

func EmailVerificationWorker() {
	for payload := range emailQueue.Out() {
		var data map[string]string
		err := json.Unmarshal([]byte(payload.(string)), &data)
		if err != nil {
			log.Printf("Error unmarshalling email payload: %s", err.Error())
			continue
		}
		email := data["email"]
		verificationLink := data["verificationLink"]

		//Send email
		if err := SendVerificationEmail(email, verificationLink); err != nil {
			log.Printf("Error sending verification email: %s", err.Error())
		}
	}
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

	var auth models.User

	if err := database.DB.Where("email = ?", data["email"]).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "account_not_found", Message: "Account not found!"}
			return c.Status(fiber.StatusNotFound).JSON(rp)
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
		"id":        auth.ID,
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

	database.DB.Model(&auth).Where("email = ?", data["email"]).Update("updated_at", time.Now()) // update the last logged in datetime

	c.Append("X-Maker-Token", signedToken)

	rp := models.ResponsePacket{Error: false, Code: "successfull", Message: "Login successfull"}
	return c.Status(fiber.StatusOK).JSON(rp)
}

// Email Verification Route
func VerifyEmail(c *fiber.Ctx) error {
	// Grab data from url
	verificationCode := c.Params("verificationCode")
	email := c.Params("email")

	if verificationCode == "" || email == "" {
		rp := models.ResponsePacket{Error: true, Code: "empty_code", Message: "Missing data in url."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	decodedEmail, err := base64.StdEncoding.DecodeString(email)
	decodedVerificationCode, err := base64.StdEncoding.DecodeString(verificationCode)
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "invalid_email", Message: "Invalid email."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var verification models.User

	if err := database.DB.Where("email = ?", decodedEmail).First(&verification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Verification code for this user does not exsist"}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	// Check if verification code has expired. If expired, send a new one.
	if uint(time.Now().Unix()) > verification.VerificationExpiry {
		if err := database.DB.Model(&verification).Where("email = ?", decodedEmail).Updates(map[string]interface{}{
			"verification_expiry": uint(time.Now().Add(time.Minute * 30).Unix())}).Error; err != nil {
			rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Verification time frame has expired, however server encountered problem while sending a new link. Please try again later."}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}

		verificationLink := fmt.Sprintf("http://localhost:8000/api/v2/verify/%s/%s", email, verificationCode)

		body := fmt.Sprintf(`{"email": "%s", "verificationLink": "%s"}`, email, verificationLink)
		emailQueue.In() <- body

		rp := models.ResponsePacket{Error: true, Code: "expired", Message: "Verfication time frame has expired. A new link has been sent to your email."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	// Check if verification code matches.
	if verification.VerificationCode != string(decodedVerificationCode) {
		rp := models.ResponsePacket{Error: true, Code: "code_mismatch", Message: "Verification code does not match."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if err := database.DB.Model(&verification).Where("email = ?", decodedEmail).Updates(map[string]interface{}{"verified": true, "verification_code": gorm.Expr("NULL")}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Unable to update verification status for the user."}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	rp := models.ResponsePacket{Error: false, Code: "verified", Message: "Verification successfull."}
	return c.Status(fiber.StatusAccepted).JSON(rp)
}

/*Update Password*/
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

/*Password Reset*/
func ResetPassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if len(data["email"]) <= 0 {
		rp := models.ResponsePacket{Error: true, Code: "missing_data", Message: "Form is missing required data!"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if !emailIsValid(data["email"]) {
		rp := models.ResponsePacket{Error: true, Code: "invalid_email", Message: "Email is not valid."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var user models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "User not found."}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	err := SendPasswordResetEmail(user.Email)
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not send email."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Email sent!"})
}

/*
 * HELPER FUNCTIONS
 */
func SendPasswordResetEmail(email string) error {

	auth := smtp.PlainAuth("", "231c63d58c7571", "15065dc065bf4c", "sandbox.smtp.mailtrap.io")

	to := []string{email}
	subject := "Subject: Password Reset\n"
	from := "maker@example.com"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<div  style="font-size:20px; font-family: Arial, serif;">
				<p>Hi there,</p>
				<p>It looks like you requested a password reset. If this was you, please click the link below to reset your password.</p>
				<p>If you did not request a password reset, please ignore this email.</p>
				<p>Reset your password here: <code style="font-weight: bold;"><button>Reset Password</button></p>
			</div>
		</html>
		`)
	msg := []byte(subject + mime + body)

	err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", auth, from, to, msg)
	return err
}

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
func SendVerificationEmail(email string, verificationLink string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "231c63d58c7571", "15065dc065bf4c", "sandbox.smtp.mailtrap.io")

	to := []string{email}
	subject := "Subject: Makers Verification Code\n"
	from := "maker@example.com"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<div  style="font-size:20px; font-family: Arial, serif;">
				<p>Clink the link below to verify your email address.</p>
				<a href="%s">Verify Email</a>
			</div>
		</html>
		`, verificationLink)
	msg := []byte(subject + mime + body)

	err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", auth, from, to, msg)
	return err
}
