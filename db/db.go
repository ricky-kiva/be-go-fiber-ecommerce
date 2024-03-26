package db

import (
	"be-go-fiber-ecommerce/entity"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func InitDb() *gorm.DB {
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
		&entity.Product{},
		&entity.Category{},
		&entity.User{},
		&entity.Cart{},
		&entity.CartItem{},
	)

	SeedData(db)
}
