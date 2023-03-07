package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	UUID           string `json:"-" gorm:"primaryKey;unique"` //uuid Generated
	Name           string `gorm:"unique"`
	UrlSafeName    string `gorm:"unique"` // Generated
	Difficulty     string `gorm:"not null"`
	Author         string `gorm:"not null"` //Use logged in email as author
	TimeToComplete string `gorm:"not null"`
	MaterialCost   string `gorm:"not null"`
	License        string `gorm:"not null"`
	Stage          string `gorm:"not null"`
	Capabilities   []Capability
	Disabilities   []Disability
	Usages         []Usage
	Images         []Image
}

type Capability struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
}

type Disability struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
}

type Usage struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DeviceID    uint
}

type File struct {
	gorm.Model
	Name             string `gorm:"not null"`
	Description      string
	DiscLocationName string
	FileSize         uint
	DeviceID         uint
}

type Image struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	DeviceID    uint
}

type Comment struct {
	gorm.Model
	Author   string
	Title    string
	Details  string
	Rating   string
	DeviceID uint
}
