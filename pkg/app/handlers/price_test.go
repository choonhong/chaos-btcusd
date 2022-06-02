package handlers

import (
	"testing"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
	"github.com/stretchr/testify/assert"
)

var timeNow = time.Now()

func TestGetBTCToUSD(t *testing.T) {
	price, err := getBTCToUSD() 
	assert.NoError(t, err)
	assert.NotEqual(t, 0, price)
}

func TestAddPrice(t *testing.T) {
	var rate model.ExchangeRate
	database.Connect()
	database.DB.Delete(&model.ExchangeRate{}, "id > 0")
	assert.NoError(t, addPrice(30287, timeNow))
	database.DB.Last(&rate)
	assert.Equal(t, 30287, rate.USD)
	assert.Equal(t, timeNow.Unix(), rate.CreatedAt.Unix())
}

func TestGetLatestPrice(t *testing.T) {
	rate := model.ExchangeRate{
		USD: 30287,
		CreatedAt: timeNow,
	}
	database.Connect()
	database.DB.Delete(&model.ExchangeRate{}, "id > 0")
	database.DB.Create(&rate)
	usd := getLatestPrice()
	assert.Equal(t, 30287, usd)
}

func TestGetPrice(t *testing.T) {
	database.Connect()
	database.DB.Delete(&model.ExchangeRate{}, "id > 0")
	rate1 := model.ExchangeRate{
		USD: 30287,
		CreatedAt: timeNow,
	}
	rate2 := model.ExchangeRate{
		USD: 30317,
		CreatedAt: timeNow.Add(time.Minute * 3),
	}
	database.DB.Create(&rate1)
	database.DB.Create(&rate2)

	t.Run("no price before", func(t *testing.T) {
		price := getPrice(timeNow.Add(-time.Hour))
		assert.Equal(t, float64(30287), price)
	})

	t.Run("no price after", func(t *testing.T) {
		price := getPrice(timeNow.Add(time.Hour))
		assert.Equal(t, float64(30317), price)
	})

	t.Run("in between", func(t *testing.T) {
		price := getPrice(timeNow.Add(time.Minute * 2))
		assert.Equal(t, float64(30307), price)
	})
}

func TestGetAveragePrice(t *testing.T) {
	database.Connect()
	database.DB.Delete(&model.ExchangeRate{}, "id > 0")
	rate1 := model.ExchangeRate{
		USD: 30287,
		CreatedAt: timeNow,
	}
	rate2 := model.ExchangeRate{
		USD: 30317,
		CreatedAt: timeNow.Add(time.Minute * 3),
	}
	database.DB.Create(&rate1)
	database.DB.Create(&rate2)

	price := getAveragePrice(timeNow, timeNow.Add(time.Minute * 3))
	assert.Equal(t, float64(30302), price)
}
