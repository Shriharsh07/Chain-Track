package config

import (
	"fmt"
	"log"

	"github.com/Shriharsh07/chaintrack/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(ConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	_ = DB.AutoMigrate(&models.Transaction{}, &models.Block{})

	fmt.Println("Database connected successfully!")
}
