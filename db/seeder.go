package db

import (
	"be-go-fiber-ecommerce/models"
	"fmt"

	"github.com/jinzhu/gorm"
)

func SeedData(db *gorm.DB) {
	category := models.Category{}
	if db.Find(&category).RecordNotFound() {
		seedCategories(db)
	}

	product := models.Product{}
	if db.Find(&product).RecordNotFound() {
		seedProducts(db)
	}
}

func seedCategories(db *gorm.DB) {
	categories := []models.Category{
		{Name: "Saltwater"},
		{Name: "Freshwater"},
		{Name: "Anemones"},
		{Name: "Rare"},
	}

	for _, category := range categories {
		db.FirstOrCreate(&category, models.Category{Name: category.Name})
	}
}

func seedProducts(db *gorm.DB) {
	products := []struct {
		Name         string
		Description  string
		Price        float64
		Stock        int
		CategoryName string
	}{
		{
			Name:         "Clownfish",
			Description:  "Vibrant orange and white stripes make this saltwater species a favourite among aquarium enthusiasts.",
			Price:        25000,
			Stock:        10,
			CategoryName: "Saltwater",
		},
		{
			Name:         "Blue Tang",
			Description:  "Known for its striking blue and yellow coloring, this saltwater fish adds a splash of colour to any tank.",
			Price:        60000,
			Stock:        5,
			CategoryName: "Saltwater",
		},
		{
			Name:         "Lionfish",
			Description:  "With its distinctive stripes and an array of venomous spines, the lionfish is both beautiful and dangerous.",
			Price:        80000,
			Stock:        8,
			CategoryName: "Saltwater",
		},
		{
			Name:         "Goldfish",
			Description:  "A classic freshwater aquarium fish, goldfish come in a variety of shapes and colours.",
			Price:        15000,
			Stock:        20,
			CategoryName: "Freshwater",
		},
		{
			Name:         "Betta Fish",
			Description:  "These small, aggressive fish are known for their vibrant colours and large, flowing fins.",
			Price:        30000,
			Stock:        12,
			CategoryName: "Freshwater",
		},
		{
			Name:         "Angelfish",
			Description:  "Tall and thin with long, flowing fins, angelfish are a popular choice for freshwater aquariums.",
			Price:        45000,
			Stock:        9,
			CategoryName: "Freshwater",
		},
		{
			Name:         "Bubble Tip Anemone",
			Description:  "The Bubble Tip Anemone is a stunning marine invertebrate prized for its vibrant coloration and unique bubble-like tips. It forms a symbiotic relationship with certain clownfish species and adds both beauty and complexity to reef aquariums.",
			Price:        15000,
			Stock:        8,
			CategoryName: "Anemones",
		},
		{
			Name:         "Carpet Anemone",
			Description:  "Carpet Anemones are large, predatory anemones found in tropical reef environments. With their wide range of colors and intricate patterns, they create a stunning centerpiece in marine aquariums. They require proper care and space due to their size and aggressive nature.",
			Price:        20000,
			Stock:        5,
			CategoryName: "Anemones",
		},
		{
			Name:         "Magnificent Sea Anemone",
			Description:  "The Magnificent Sea Anemone is renowned for its majestic appearance and striking coloration. Found in the warm waters of the Indo-Pacific region, it forms a symbiotic relationship with various clownfish species. Its long, flowing tentacles sway gracefully in the currents, making it a captivating addition to reef aquariums.",
			Price:        30000,
			Stock:        3,
			CategoryName: "Anemones",
		},
		{
			Name:         "Arowana",
			Description:  "Often called the 'dragon fish', arowanas are large, predatory fish valued for their unique appearance.",
			Price:        250000,
			Stock:        3,
			CategoryName: "Rare",
		},
		{
			Name:         "Flowerhorn Cichlid",
			Description:  "Known for its vivid colours and the large nuchal hump on its head, the Flowerhorn Cichlid is a prized rare fish.",
			Price:        120000,
			Stock:        4,
			CategoryName: "Rare",
		},
		{
			Name:         "Peppermint Angelfish",
			Description:  "One of the rarest aquarium fish, the Peppermint Angelfish is renowned for its striking red and white stripes.",
			Price:        300000,
			Stock:        2,
			CategoryName: "Rare",
		},
	}

	for _, p_temp := range products {
		var category models.Category

		if err := db.Where("name = ?", p_temp.CategoryName).First(&category).Error; err != nil {
			fmt.Println("Error finding category for product:", p_temp.Name)
			continue
		}

		product := models.Product{}

		db.FirstOrCreate(&product, models.Product{
			Name:        p_temp.Name,
			Description: p_temp.Description,
			Price:       p_temp.Price,
			Stock:       p_temp.Stock,
			CategoryID:  category.ID,
		})
	}
}
