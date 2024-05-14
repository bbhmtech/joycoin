package model

import "time"

type Transaction struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FromAccountID uint
	ToAccountID   uint
	Message       string
	CentAmount    int64
}
