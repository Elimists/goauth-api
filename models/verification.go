package models

type Verification struct {
	Code    string
	Email   string `gorm:"unique"`
	Expires uint
}
