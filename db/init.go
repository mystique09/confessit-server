package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected.")
	return conn
}
