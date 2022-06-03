package handlers

import (
	"testing"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
	"github.com/stretchr/testify/assert"
)

var timeNow = time.Now()

func TestFetchPrice(t *testing.T) {
	var rate model.ExchangeRate
	database.Connect()
	database.DB.Delete(&model.ExchangeRate{}, "id > 0")
	FetchPrice() 
	database.DB.Last(&rate)
	assert.NotEqual(t, 0, rate.USD)
	assert.NotEqual(t, 0, lastPrice)
}

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
