package main

import (
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
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	// fetch exchange rate once per minute
	handlers.FetchPrice()
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
	r.Get("/", handlers.GetPrice)

	http.ListenAndServe(":80", r)
}
