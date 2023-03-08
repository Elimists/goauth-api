package models

type Request struct {
	ID              string `gorm:"primaryKey"`
	RequestedDevice string // reference device
	RequestedBy     string // reference user
	Title           string
	Details         string
	RequestedDate   uint
	UpdatedDate     uint
}
