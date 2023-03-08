package models

type User struct {
	Email        string `gorm:"unique"`
	Name         string
	Education    string
	Occupation   string
	Organization string
	Bio          string
	Picture      string
}
