package models

type Address struct {
	Id      uint `gorm:"unique"`
	Email   string
	Name    string
	Street  string
	City    string
	State   string
	Country string
	Postal  string
	Active  bool
}
