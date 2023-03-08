package database

import (
	"fmt"
	"os"

	"github.com/Elimists/go-app/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	/*
		user := os.Getenv("MYSQL_USER")
		pass := os.Getenv("MYSQL_PASS")
		connection := fmt.Sprintf("%s:%s@/devdb?parseTime=true", user, pass)
	*/

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	stringConn := fmt.Sprintf("%s:%s@/%s?parseTime=true", user, password, database)

	db_conn, err := gorm.Open(mysql.Open(stringConn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Could not connect to database")
	}

	DB = db_conn

	db_conn.AutoMigrate(&models.Auth{})
	db_conn.AutoMigrate(&models.Verification{})
	db_conn.AutoMigrate(&models.User{})
	db_conn.AutoMigrate(&models.Device{})
	db_conn.AutoMigrate(&models.Capability{})
	db_conn.AutoMigrate(&models.Device{})
	db_conn.AutoMigrate(&models.Disability{})
	db_conn.AutoMigrate(&models.Usage{})
	db_conn.AutoMigrate(&models.File{})
	db_conn.AutoMigrate(&models.Image{})
	db_conn.AutoMigrate(&models.Comment{})
}
