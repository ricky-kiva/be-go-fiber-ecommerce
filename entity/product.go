package entity

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    Category
	CategoryID  uint
}
