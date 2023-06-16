package models

import (
	"time"
)

// CustomModel is used instead of gorm.Model to avoid the DeletedAt fields.
type CustomModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	CustomModel
	Email            string           `json:"email" gorm:"unique;primary_key"` // The email address of the user. This is the primary key.
	Password         []byte           `json:"-"`
	Privilege        int8             `json:"privilege"` // 1: Admin, 2: Manager, 3: Coordinator, 4: Moderator, 9: General user
	Verified         bool             `json:"-"`
	UserVerification UserVerification `json:"userVerification" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
	UserDetails      UserDetails      `json:"userDetails" gorm:"constraint:OnDelete:CASCADE;foreignkey:UserID"` // One to one relationship with the user details. Delete the user details if the user is deleted.
}

type UserVerification struct {
	CustomModel
	UserID             uint   `json:"-"` // The user ID of the user this verification belongs to.
	VerificationCode   string `json:"-"`
	VerificationExpiry uint   `json:"-"`
}

type UserDetails struct {
	CustomModel
	UserID         uint               `json:"-"`                                         // The user ID of the user this verification belongs to.
	UserEmail      string             `json:"userEmail" gorm:"unique;type:varchar(255)"` // The email address of the user. This is the primary key.
	FirstName      string             `json:"firstName"`
	LastName       string             `json:"lastName"`
	Bio            string             `json:"bio"`
	Interests      string             `json:"interests"`
	Organization   string             `json:"organization"` // The organization the user is affiliated with.
	Occupation     string             `json:"occupation"`
	Education      string             `json:"education"` // The level of education the user has achieved
	ProfilePicture UserProfilePicture `json:"profilePicture" gorm:"constraint:OnDelete:CASCADE;foreignkey:UserDetailsID"`
	Addresses      []UserAddress      `json:"addresses;" gorm:"constraint:OnDelete:CASCADE;foreignkey:UserDetailsID"`
}

type UserAddress struct {
	CustomModel
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	State         string `json:"state"`
	ZipCode       string `json:"zipCode"`
	Country       string `json:"country"`
	IsActive      bool   `json:"isActive"`
	UserDetailsID uint   `json:"userDetailsID"`
}
type UserProfilePicture struct {
	CustomModel
	UrlSafeName   string `json:"urlSafeName"`
	UserImage     []byte `json:"userImage"`
	ImageAltText  string `json:"imageAltText"`
	UserDetailsID uint   `json:"userDetailsID"`
}
