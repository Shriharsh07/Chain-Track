package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Transaction struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36) primaryKey"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Amount    float64   `json:"amount"`
	IsMined   bool      `gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
