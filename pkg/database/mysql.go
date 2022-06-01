package database

import (
	"github.com/chaos-btcusd/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error
	dsn := "root:root@tcp(btcusd-db:3306)/pricing?charset=utf8mb4&parseTime=True&loc=Local"
  DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&model.ExchangeRate{})
}
