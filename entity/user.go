package entity

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(255);unique_index"`
	Password string `gorm:"size:255"`
	Cart     Cart
}
