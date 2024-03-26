package entity

import "github.com/jinzhu/gorm"

type Cart struct {
	gorm.Model
	Items  []CartItem
	UserID uint `gorm:"unique"`
}
