package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
	"github.com/chaos-btcusd/pkg/utils"
	"gorm.io/gorm"
)

// GetPrice returns BTC-USD price
func GetPrice(w http.ResponseWriter, r *http.Request) {
	var price float64
	var err error
	layout := "2006-01-02T15:04:05Z"
	timestamp := r.URL.Query().Get("timestamp")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	switch {
	case timestamp != "":
		// given timestamp, get price
		requestedTime, err1 := time.Parse(layout, timestamp)
		if err1 != nil {
			fmt.Println(err1)
			utils.ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}
		price, err = getPriceAtTime("bitcoin", requestedTime)

	case from != "" && to != "":
		// given time range, get average price
		fromTime, err1 := time.Parse(layout, from)
		toTime, err2 := time.Parse(layout, to)
		if err1 != nil || err2 != nil {
			utils.ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}
		price, err = getAveragePrice("bitcoin", fromTime, toTime)

	default:
		// get last price
		price, err = getLatestPrice("bitcoin")
	}

	// check error
	if err != nil {
		fmt.Println("failed to get price from database, error:", err)
		utils.ResponseJSON(w, http.StatusInternalServerError, nil)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, price)
}

// get latest price data
func getLatestPrice(coin string) (float64, error) {
	var rate model.ExchangeRate
	if err := database.DB.Last(&rate, "coin = ?", coin).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return rate.USD, nil
}

// get price at a given timestamp
func getPriceAtTime(coin string, requestedTime time.Time) (float64, error) {
	var rate1, rate2 model.ExchangeRate
	if err := database.DB.Last(&rate1, "coin = ? and created_at <= ?", coin, requestedTime).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if rate1.CreatedAt.Unix() == requestedTime.Unix() {
		return rate1.USD, nil
	}
	if err := database.DB.First(&rate2, "coin = ? and created_at > ?", coin, requestedTime).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	// no price before requested time
	if rate1.ID == 0 {
		return rate2.USD, nil
	}

	// no price after requested time
	if rate2.ID == 0 {
		return rate1.USD, nil
	}

	// calculate time weighted average price
	timeGap1 := requestedTime.Sub(rate1.CreatedAt).Seconds()
	timeGap2 := rate2.CreatedAt.Sub(rate1.CreatedAt).Seconds()
	return rate1.USD + (rate2.USD-rate1.USD)*timeGap1/timeGap2, nil
}

// get average price in a given time range
func getAveragePrice(coin string, from, to time.Time) (float64, error) {
	var result float64
	if err := database.DB.Table("exchange_rates").Where("coin = ? and created_at >= ? AND created_at <= ?", coin, from, to).Select("AVG(usd)").Row().Scan(&result); err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return result, nil
}
