package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Title    string `json:"title"`
	Review   string `json:"review"`
	Rating   uint8  `json:"rating"`
	UserID   uint   `json:"userID"`
	User     User   `gorm:"constraint:OnDelete:SetNull;"` // If the user is deleted, set the review's user to null.
	DeviceID uint   `json:"deviceID" gorm:"unique"`
	Device   Device `gorm:"constraint:OnDelete:CASCADE;"` // If the device is deleted, delete the review.
}
