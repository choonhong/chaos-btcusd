package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
	"github.com/chaos-btcusd/pkg/utils"
)

var (
	lastPrice int
	lastTime 	time.Time
)

func FetchPrice() {
	usd, err := getBTCToUSD()
	if err != nil {
		fmt.Println("failed to get price, error: ", err)
		return
	}

	lastPrice = usd
	lastTime = time.Now().Round(time.Second)
	if err := addPrice(usd, lastTime); err != nil {
		fmt.Println("failed to insert price into database, error: ", err)
		return
	}
	fmt.Println(lastPrice, lastTime)
}

func getBTCToUSD() (int, error) {
	res, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if err != nil {
		return 0, err
	}

	var price map[string]map[string]int
	err = json.NewDecoder(res.Body).Decode(&price)
	if err != nil {
		return 0, err
	}

	return price["bitcoin"]["usd"], nil
}

func addPrice(usd int, createdAt time.Time) error {
	rate := model.ExchangeRate{
		USD: 			 usd,
		CreatedAt: createdAt,
	}
	return database.DB.Create(&rate).Error
}

func GetPrice(w http.ResponseWriter, r *http.Request) {
	layout := "2006-01-02T15:04:05Z"

	// get price at a given timestamp
	timestamp := r.URL.Query().Get("timestamp")
	if timestamp != "" {
		requestedTime, err := time.Parse(layout, timestamp)
		if err != nil {
			fmt.Println(err)
			utils.ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}
		utils.ResponseJSON(w, http.StatusOK, getPrice(requestedTime))
		return
	}

	// get average price in a given time range
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from != "" && to != "" {
		fromTime, err1 := time.Parse(layout, from)
		toTime, err2 := time.Parse(layout, to)
		if err1 != nil || err2 != nil {
			fmt.Println(err1)
			fmt.Println(err2)
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

func getLatestPrice() int {
	var rate model.ExchangeRate
	database.DB.Last(&rate)
	return rate.USD
}

func getPrice(requestedTime time.Time) float64 {
	var rate1, rate2 model.ExchangeRate
	database.DB.Last(&rate1, "created_at <= ?", requestedTime)
	if rate1.CreatedAt == requestedTime {
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

func getAveragePrice(from, to time.Time) float64 {
	var rate1, rate2 model.ExchangeRate
	database.DB.First(&rate1, "created_at >= ?", from)
	database.DB.Last(&rate2, "created_at <= ?", to)

	var result float64
	row := database.DB.Table("exchange_rates").Where("id between ? AND ?", rate1.ID, rate2.ID).Select("AVG(usd)").Row()
	row.Scan(&result)
	return result
}
