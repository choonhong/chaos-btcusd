package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chaos-btcusd/pkg/app/handlers"
	"github.com/chaos-btcusd/pkg/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func init() {
	// connect to db
	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect mysql successfully")

	// fetch exchange rate once per minute
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			handlers.FetchPrice()
		}
	}()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/price", handlers.GetPrice)

	http.ListenAndServe(":80", r)
}
