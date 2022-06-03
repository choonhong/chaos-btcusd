package database

import (
	"fmt"

	"github.com/chaos-btcusd/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	fmt.Println("connect to database successfully")
	return DB.AutoMigrate(&model.ExchangeRate{})
}
