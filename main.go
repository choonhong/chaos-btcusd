package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	handlers.FetchPrices()
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			handlers.FetchPrices()
		}
	}()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{coin}", handlers.GetPrice)

	if http.ListenAndServe(":80", r) != nil {
		fmt.Println("serve on :80 error")
		http.ListenAndServe(":"+os.Getenv("PORT"), r)
	}
}
