package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Title    string `json:"title"`
	Review   string `json:"review"`
	UserID   uint   `json:"userID"`
	DeviceID uint   `json:"deviceID" gorm:"unique"`
}
