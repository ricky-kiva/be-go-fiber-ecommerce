package models

import "github.com/jinzhu/gorm"

type CartItem struct {
	gorm.Model
	Quantity  int
	Product   Product
	ProductID uint
	CartID    uint
}
