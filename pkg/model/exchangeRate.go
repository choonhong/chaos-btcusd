package model

import (
	"time"
)

type ExchangeRate struct {
	ID   			uint
	USD  		  int
	CreatedAt time.Time `gorm:"index"`
}
