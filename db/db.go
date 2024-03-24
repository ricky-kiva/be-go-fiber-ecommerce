package db

import (
	"be-go-fiber-ecommerce/models"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDb() *gorm.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error load .env")
	}

	conn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", conn)

	if err != nil {
		log.Fatal("Error connecting to DB")
	}

	migrateDB(db)

	return db
}

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.User{},
	)

	SeedData(db)
}
