package models

import (
	"gorm.io/gorm"
)

type Block struct {
	gorm.Model
	Transactions []Transaction `gorm:"-"`
	PreviousHash string
	Hash         string
	Nonce        int
}
