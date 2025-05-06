package models

import "time"

type Block struct {
	ID           uint `gorm:"primaryKey"`
	Timestamp    time.Time
	Transactions []Transaction `gorm:"-"`
	PreviousHash string
	Hash         string
}
