package models

/**
* The Auth model describes the fields related to user
* authentication and authorization.
 */
type Auth struct {
	Id           uint
	Email        string `gorm:"unique"`
	Password     []byte `json:"-"`
	Verified     bool
	Privilege    int8 // 1: Admin, 2: Manager, 3: Coordinator, 4: Moderator, 9: General user
	RegisteredOn uint
	LastLoggedIn uint
}
