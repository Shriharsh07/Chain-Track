package models

import (
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

type Transaction struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Sender    string
	Receiver  string
	Amount    float64
	CreatedAt time.Time
}
