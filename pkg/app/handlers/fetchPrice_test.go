package handlers

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
	"github.com/stretchr/testify/assert"
)

var timeNow = time.Now()

func TestFetchPrice(t *testing.T) {
	database.Connect()
	defer os.Remove("./gorm.db")
	FetchPrices()
	for _, coin := range supportedCoins {
		fmt.Println(coin)
		var rate model.ExchangeRate
		database.DB.Last(&rate, "coin = ?", coin)
		assert.Equal(t, coin, rate.Coin)
		assert.NotEqual(t, 0, rate.USD)
	}
}

func TestGetCoinToUSD(t *testing.T) {
	price, err := getCoinToUSD("bitcoin")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, price)
}

func TestAddPrice(t *testing.T) {
	t.Run("db error", func(t *testing.T) {
		assert.Error(t, addPrice("bitcoin", 30287, timeNow))
	})

	t.Run("normal", func(t *testing.T) {
		var rate model.ExchangeRate
		database.Connect()
		defer os.Remove("./gorm.db")
		assert.NoError(t, addPrice("bitcoin", 30287, timeNow))
		database.DB.Last(&rate)
		assert.Equal(t, "bitcoin", rate.Coin)
		assert.Equal(t, float64(30287), rate.USD)
		assert.Equal(t, timeNow.Unix(), rate.CreatedAt.Unix())
	})

}
