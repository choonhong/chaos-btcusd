package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
	"github.com/stretchr/testify/assert"
)

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

func TestGetPriceWithTime(t *testing.T) {
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
		price := getPriceWithTime(timeNow.Add(-time.Hour))
		assert.Equal(t, float64(30287), price)
	})

	t.Run("no price after", func(t *testing.T) {
		price := getPriceWithTime(timeNow.Add(time.Hour))
		assert.Equal(t, float64(30317), price)
	})

	t.Run("have price at timestamp", func(t *testing.T) {
		price := getPriceWithTime(timeNow)
		assert.Equal(t, float64(30287), price)
	})

	t.Run("in between", func(t *testing.T) {
		price := getPriceWithTime(timeNow.Add(time.Minute * 2))
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

// TestGetPrice assume all above tests passes
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


	t.Run("get last price", func(t *testing.T) {
		req := httptest.NewRequest("GET", "localhost:80/price", nil)
		w := httptest.NewRecorder()
		GetPrice(w, req)

		var price map[string]int
		res := w.Result()
		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&price))
		assert.Equal(t, 30317, price["data"])
	})

	t.Run("get last price from db", func(t *testing.T) {
		req := httptest.NewRequest("GET", "localhost:80/price", nil)
		w := httptest.NewRecorder()
		lastPrice = 20317
		GetPrice(w, req)

		var price map[string]int
		res := w.Result()
		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&price))
		assert.Equal(t, 20317, price["data"])
	})

	t.Run("bad timestamp", func(t *testing.T) {
		req := httptest.NewRequest("GET", "localhost:80/price?timestamp=2022-06-01T18:39:", nil)
		w := httptest.NewRecorder()
		GetPrice(w, req)

		res := w.Result()
		assert.Equal(t, 400, res.StatusCode)
	})

	t.Run("get price at timestamp", func(t *testing.T) {
		req := httptest.NewRequest("GET", "localhost:80/price?timestamp=2022-06-01T18:39:47Z", nil)
		w := httptest.NewRecorder()
		GetPrice(w, req)

		var price map[string]int
		res := w.Result()
		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&price))
		assert.Equal(t, 30287, price["data"])
	})

	t.Run("bad time range", func(t *testing.T) {
		req := httptest.NewRequest("GET", "localhost:80/price?from=2022-06-01T18:39:04Z&to=2023-06-01T18:47:", nil)
		w := httptest.NewRecorder()
		GetPrice(w, req)

		res := w.Result()
		assert.Equal(t, 400, res.StatusCode)
	})

	t.Run("get average price", func(t *testing.T) {
		req := httptest.NewRequest("GET", "localhost:80/price?from=2022-06-01T18:39:04Z&to=2023-06-01T18:47:47Z", nil)
		w := httptest.NewRecorder()
		GetPrice(w, req)

		var price map[string]int
		res := w.Result()
		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&price))
		assert.Equal(t, 30302, price["data"])
	})
}
