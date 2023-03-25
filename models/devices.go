package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name           string `gorm:"unique"`
	UrlSafeName    string `gorm:"unique"` // Generated
	Difficulty     string `gorm:"not null"`
	TimeToComplete string `gorm:"not null"`
	MaterialCost   string `gorm:"not null"`
	License        string `gorm:"not null"`
	Stage          string `gorm:"not null"`
	Capabilities   []DeviceCapability
	Disabilities   []DeviceDisability
	Usages         []DeviceUsage
	Images         []DeviceImage
	Reviews        []Review // List of reviews users have submitted for this device
	UserID         uint     `json:"userID"`
	User           User     `gorm:"constraint:OnDelete:SetNull;"` // If the user is deleted, set the device's user to null.
}

type DeviceCapability struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
	Device      Device `gorm:"constraint:OnDelete:CASCADE;"` // If the device is deleted, delete the device's capabilities.
}

type DeviceDisability struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
	Device      Device `gorm:"constraint:OnDelete:CASCADE;"` // If the device is deleted, delete the device's disabilities.
}

type DeviceUsage struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
	Device      Device `gorm:"constraint:OnDelete:CASCADE;"` // If the device is deleted, delete the device's usages.
}

type DeviceFile struct {
	gorm.Model
	Name             string `gorm:"not null"`
	Description      string
	DiscLocationName string
	FileSize         uint
	DeviceID         uint
	Device           Device `gorm:"constraint:OnDelete:CASCADE;"` // If the device is deleted, delete the device's files.
}

type DeviceImage struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	DeviceID    uint
	Device      Device `gorm:"constraint:OnDelete:CASCADE;"` // If the device is deleted, delete the device's images.
}
