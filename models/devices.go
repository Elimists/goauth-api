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
	Reviews        []Review
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
	DeviceID    uint
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
