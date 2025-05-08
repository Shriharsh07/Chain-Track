package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	maxRetries := 10
	retryInterval := 5 * time.Second

	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")

	// Build connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db_user,
		db_pass,
		db_host,
		db_port,
		db_name,
	)

	// Retry connection logic
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Attempt %d: Failed to connect to database: %v", i+1, err)
		if i < maxRetries-1 {
			log.Printf("Waiting %v before retrying...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		log.Fatal("Failed to connect to database after retries: ", err)
	}

	// Test connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	fmt.Println("Database connected successfully!")
}
