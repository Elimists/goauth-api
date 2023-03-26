package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name           string             `gorm:"unique"`
	UrlSafeName    string             `gorm:"unique"` // Generated
	Difficulty     string             `gorm:"not null"`
	TimeToComplete string             `gorm:"not null"`
	MaterialCost   string             `gorm:"not null"`
	License        string             `gorm:"not null"`
	Stage          string             `gorm:"not null"`
	Capabilities   []DeviceCapability `gorm:"constraint:OnDelete:CASCADE;"` // If the device is deleted, delete the capabilities for this device.
	Disabilities   []DeviceDisability `gorm:"constraint:OnDelete:CASCADE;"`
	Usages         []DeviceUsage      `gorm:"constraint:OnDelete:CASCADE;"`
	Images         []DeviceImage      `gorm:"constraint:OnDelete:CASCADE;"`
	Reviews        []Review           `gorm:"constraint:OnDelete:SET NULL;"` // If the device is deleted, set the device ID for the reivew to null. The user is responsible for deleting the review.
	UserDetailsID  uint               `json:"userID"`
}

type DeviceCapability struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
}

type DeviceDisability struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
}

type DeviceUsage struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint // References the device in the Device table.
}

type DeviceFile struct {
	gorm.Model
	Name             string `gorm:"not null"`
	Description      string
	DiscLocationName string
	FileSize         uint
	DeviceID         uint
}

type DeviceImage struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	DeviceID    uint
}
