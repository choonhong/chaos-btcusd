package database

import (
	"fmt"

	"github.com/chaos-btcusd/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	fmt.Println("connect database successfully")
	return DB.AutoMigrate(&model.ExchangeRate{})
}
