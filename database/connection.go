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
	db_conn.AutoMigrate(&models.Review{})
	db_conn.AutoMigrate(&models.Device{}, &models.DeviceCapability{}, &models.DeviceDisability{}, &models.DeviceUsage{}, &models.DeviceFile{}, &models.DeviceImage{})
}
