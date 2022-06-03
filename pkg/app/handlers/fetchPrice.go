package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
)

var (
	lastPrice int
	lastTime 	time.Time
)

// FetchPrice gets BTC-USD and stores price in database
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

// get BTC-USD price
func getBTCToUSD() (int, error) {
	res, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if err != nil {
		return 0, err
	}

	var price map[string]map[string]int
	if err := json.NewDecoder(res.Body).Decode(&price); err != nil {
		return 0, err
	}

	return price["bitcoin"]["usd"], nil
}

// insert price into database
func addPrice(usd int, createdAt time.Time) error {
	rate := model.ExchangeRate{
		USD: 			 usd,
		CreatedAt: createdAt,
	}
	return database.DB.Create(&rate).Error
}
