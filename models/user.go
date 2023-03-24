package models

import "gorm.io/gorm"

type UserAuth struct {
	gorm.Model
	Email              string `gorm:"unique"`
	Password           []byte `json:"-"`
	Privilege          int8   // 1: Admin, 2: Manager, 3: Coordinator, 4: Moderator, 9: General user
	Verified           bool
	VerificationCode   string
	VerificationExpiry uint
	UserID             uint
	User               User
}

type User struct {
	gorm.Model
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Bio          string        `json:"bio"`
	Interests    string        `json:"interests"`
	Organization string        `json:"organization"`
	Occupation   string        `json:"occupation"`
	Education    string        `json:"education"`
	Devices      []Device      `json:"devices"` // List of devices that the user has submitted.
	Addresses    []UserAddress `json:"addresses"`
	Requests     []UserRequest `json:"requests"` // List of devices the user has requested.
	Makes        []UserMake    `json:"makes"`    // List of devices the user has helped make.
	Ideas        []UserIdea    `json:"ideas"`    // List of ideas and or suggestions the user has submitted.
	Reviews      []Review      `json:"reviews"`  // List of reviews the user has submitted for various devices.
}

type UserAddress struct {
	gorm.Model
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	State         string `json:"state"`
	ZipCode       string `json:"zipCode"`
	Country       string `json:"country"`
	IsActive      bool   `json:"isActive"`
	UserID        uint   `json:"userID"`
}

type UserRequest struct {
	gorm.Model
	UserID   uint   `json:"userID"`
	DeviceID uint   `json:"deviceID"`
	Device   Device `json:"device"`
}

type UserMake struct {
	gorm.Model
	UserID   uint   `json:"userID"`
	DeviceID uint   `json:"deviceID"`
	Device   Device `json:"device"`
}

type UserIdea struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint   `json:"userID"`
}

type UserProfilePicture struct {
	gorm.Model
	UrlSafeName  string `json:"urlSafeName"`
	UserImage    []byte `json:"userImage"`
	ImageAltText string `json:"imageAltText"`
	UserID       uint   `json:"userID"`
}
