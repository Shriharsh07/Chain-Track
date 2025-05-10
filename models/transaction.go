package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Transaction struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Sender    string    `json:"sender" validate:"required,email"`
	Receiver  string    `json:"receiver" validate:"required,email"`
	Amount    float64   `json:"amount" validate:"required,min=0.01"`
	IsMined   bool      `gorm:"default:false"`
	BlockID   uint      `json:"block_id"`
	CreatedAt time.Time `json:"created_at"`
}
