package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Title         string `json:"title"`
	Review        string `json:"review"`
	Rating        uint8  `json:"rating"`
	UserDetailsID uint   `json:"userID"`
	DeviceID      uint   `json:"deviceID" gorm:"unique"`
}
