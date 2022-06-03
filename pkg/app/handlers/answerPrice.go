package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
	"github.com/chaos-btcusd/pkg/utils"
)

// GetPrice returns BTC-USD price
func GetPrice(w http.ResponseWriter, r *http.Request) {
	layout := "2006-01-02T15:04:05Z"

	// given timestamp, get price
	timestamp := r.URL.Query().Get("timestamp")
	if timestamp != "" {
		requestedTime, err := time.Parse(layout, timestamp)
		if err != nil {
			fmt.Println(err)
			utils.ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}
		utils.ResponseJSON(w, http.StatusOK, getPriceWithTime(requestedTime))
		return
	}

	// given time range, get average price
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from != "" && to != "" {
		fromTime, err1 := time.Parse(layout, from)
		toTime, err2 := time.Parse(layout, to)
		if err1 != nil || err2 != nil {
			utils.ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}
		utils.ResponseJSON(w, http.StatusOK, getAveragePrice(fromTime, toTime))
		return
	}

	// get last price
	if lastPrice != 0 {
		utils.ResponseJSON(w, http.StatusOK, lastPrice)
		return
	}
	utils.ResponseJSON(w, http.StatusOK, getLatestPrice())
}

// get latest price data
func getLatestPrice() int {
	var rate model.ExchangeRate
	database.DB.Last(&rate)
	return rate.USD
}

// get price at a given timestamp
func getPriceWithTime(requestedTime time.Time) float64 {
	var rate1, rate2 model.ExchangeRate
	database.DB.Last(&rate1, "created_at <= ?", requestedTime)
	if rate1.CreatedAt.Unix() == requestedTime.Unix() {
		return float64(rate1.USD)
	}
	database.DB.First(&rate2, "created_at > ?", requestedTime)

	// no price before requested time
	if rate1.ID == 0 {
		return float64(rate2.USD)
	}

	// no price after requested time
	if rate2.ID == 0 {
		return float64(rate1.USD)
	}

	// calculate time weighted average price
	timeGap1 := requestedTime.Sub(rate1.CreatedAt).Seconds()
	timeGap2 := rate2.CreatedAt.Sub(rate1.CreatedAt).Seconds()
	return float64(rate1.USD) + float64(rate2.USD - rate1.USD) * timeGap1 / timeGap2
}

// get average price in a given time range
func getAveragePrice(from, to time.Time) float64 {
	var rate1, rate2 model.ExchangeRate
	database.DB.First(&rate1, "created_at >= ?", from)
	database.DB.Last(&rate2, "created_at <= ?", to)

	var result float64
	database.DB.Table("exchange_rates").Where("id between ? AND ?", rate1.ID, rate2.ID).Select("AVG(usd)").Row().Scan(&result)
	return result
}