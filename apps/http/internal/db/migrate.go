package main

import (
	"fmt"
	"log"
	"os"

	"meet/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("../.env")
	fmt.Println(err, "env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database := os.Getenv("DB_URL")
	fmt.Println(database, "database url")
	db, err := gorm.Open(postgres.Open(database), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Meeting{},
		&models.MeetingParticipant{},
		&models.Message{},
	); err != nil {
		log.Fatal("migration failed:", err)
	}
	fmt.Println("Schema migrated Successfully")
}
