package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"type:VARCHAR(255)"`
	Email string `gorm:"type:VARCHAR(255)"`
	Image string `gorm:"type:TEXT"`
}
