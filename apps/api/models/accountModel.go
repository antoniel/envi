package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	CompoundID         string `gorm:"type:VARCHAR(255)"`
	UserID             int
	PoviderType        string `gorm:"type:VARCHAR(255)"`
	PoviderID          string `gorm:"type:VARCHAR(255)"`
	PoviderAccountID   string `gorm:"type:VARCHAR(255)"`
	RefreshToken       string `gorm:"type:TEXT"`
	AccessToken        string `gorm:"type:TEXT"`
	AccessTokenExpires time.Time
}

func CreateAccount(db *gorm.DB, account *Account) error {
	return db.Create(account).Error
}
