package models

import (
	"time"

	"gorm.io/gorm"
)

type Block struct {
	gorm.Model
	Transactions []Transaction `gorm:"-"`
	PreviousHash string
	Hash         string
	Nonce        int
}

type BlockResponse struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`

	Transactions []Transaction `gorm:"-" json:"transactions"`
	PreviousHash string        `json:"previous_hash"`
	Hash         string        `json:"hash"`
	Nonce        int           `json:"nonce"`
}
