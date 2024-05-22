package model

import "gorm.io/gorm"

func AutoMigration(db *gorm.DB) {
	db.AutoMigrate(&Jumper{}, &Account{}, &ShortenLink{}, &Transaction{}, &QuickAction{})
}
