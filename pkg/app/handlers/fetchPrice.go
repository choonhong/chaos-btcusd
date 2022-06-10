package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-btcusd/pkg/database"
	"github.com/chaos-btcusd/pkg/model"
)

var supportedCoins = []string{"bitcoin", "ethereum"}

// FetchPrices gets BTC-USD and stores price in database
func FetchPrices() {
	for _, coin := range supportedCoins {
		timeNow := time.Now().Round(time.Second)
		usd, err := getCoinToUSD(coin)
		if err != nil || usd == 0 {
			fmt.Println("failed to get price from 3rd party API, error: ", err)
			continue
		}

		if err := addPrice(coin, usd, timeNow); err != nil {
			fmt.Println("failed to insert price into database, error: ", err)
			continue
		}

		fmt.Println(coin, usd, timeNow)
	}
}

// get price
func getCoinToUSD(coin string) (float64, error) {
	res, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=" + coin + "&vs_currencies=usd")
	if err != nil {
		return 0, err
	}

	var price map[string]map[string]float64
	if err := json.NewDecoder(res.Body).Decode(&price); err != nil {
		return 0, err
	}

	return price[coin]["usd"], nil
}

// insert price into database
func addPrice(coin string, usd float64, createdAt time.Time) error {
	rate := model.ExchangeRate{
		Coin:      coin,
		USD:       usd,
		CreatedAt: createdAt,
	}
	return database.DB.Create(&rate).Error
}
