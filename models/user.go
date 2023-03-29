package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email              string      `json:"email" gorm:"unique"`
	Password           []byte      `json:"-"`
	Privilege          int8        `json:"privilege"` // 1: Admin, 2: Manager, 3: Coordinator, 4: Moderator, 9: General user
	Verified           bool        `json:"-"`
	VerificationCode   string      `json:"-"`
	VerificationExpiry uint        `json:"-"`
	UserDetails        UserDetails `json:"userDetails" gorm:"constraint:OnDelete:CASCADE;"` // One to one relationship with the user details. Delete the user details if the user is deleted.
}

type UserDetails struct {
	gorm.Model
	FirstName         string            `json:"firstName"`
	LastName          string            `json:"lastName"`
	Bio               string            `json:"bio"`
	Interests         string            `json:"interests"`
	Organization      string            `json:"organization"` // The organization the user is affiliated with.
	Occupation        string            `json:"occupation"`
	Education         string            `json:"education"`                                     // The level of education the user has achieved
	Devices           []Device          `json:"devices;" gorm:"constraint:OnDelete:SET NULL;"` // List of devices that the user has submitted.
	Addresses         []UserAddress     `json:"addresses;" gorm:"constraint:OnDelete:CASCADE;"`
	Requests          []UserRequest     `json:"requests;" gorm:"constraint:OnDelete:CASCADE"` // List of devices the user has requested.
	Makes             []UserMake        `json:"makes;" gorm:"constraint:OnDelete:CASCADE"`    // List of devices the user has helped make.
	Ideas             []UserIdea        `json:"ideas;" gorm:"constraint:OnDelete:CASCADE"`    // List of ideas and or suggestions the user has submitted.
	Reviews           []Review          `json:"reviews;" gorm:"constraint:OnDelete:CASCADE"`  // List of reviews the user has submitted for various devices.
	UserID            uint              `json:"userID" gorm:"unique"`
	UserAccountStatus UserAccountStatus `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}

type UserAddress struct {
	gorm.Model
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	State         string `json:"state"`
	ZipCode       string `json:"zipCode"`
	Country       string `json:"country"`
	IsActive      bool   `json:"isActive"`
	UserDetailsID uint   `json:"userDetailsID"`
}

type UserRequest struct {
	gorm.Model
	DeviceID      uint   `json:"deviceID"`
	Device        Device `json:"device"`
	UserDetailsID uint   `json:"userDetailsID"`
}

type UserMake struct {
	gorm.Model
	DeviceID      uint   `json:"deviceID"`
	Device        Device `json:"device"`
	UserDetailsID uint   `json:"userDetailsID"`
}

type UserIdea struct {
	gorm.Model
	Title         string `json:"title"`
	Content       string `json:"content"`
	UserDetailsID uint   `json:"userDetailsID"`
}

type UserProfilePicture struct {
	gorm.Model
	UrlSafeName   string `json:"urlSafeName"`
	UserImage     []byte `json:"userImage"`
	ImageAltText  string `json:"imageAltText"`
	UserDetailsID uint   `json:"userDetailsID"`
}

type UserAccountStatus struct {
	gorm.Model
	IsSuspended   bool   `json:"isSuspended" gorm:"required"`
	SuspensionMsg string `json:"suspensionMsg" gorm:"required"`
	UserDetailsID uint   `json:"userDetailsID"`
}
