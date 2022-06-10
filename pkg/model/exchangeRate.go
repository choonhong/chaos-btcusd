package model

import (
	"time"
)

type ExchangeRate struct {
	ID        uint
	Coin      string
	USD       float64
	CreatedAt time.Time `gorm:"index"`
}
