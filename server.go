package main

import (
	"genesisbankly/exchange/db"
	"genesisbankly/exchange/routes"
	"genesisbankly/exchange/routes/adapters"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func setDefaultRoutes(router *chi.Mux, connection *gorm.DB) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		adapters.RespondWithJson(w, 200, struct{}{})
	})
	router.Get("/errors", func(w http.ResponseWriter, r *http.Request) {
		logs, err := db.ListLogsInsideDb(connection)
		if err != nil {
			adapters.RespondWithError(w, connection, 500, err, r)
		} else {
			adapters.RespondWithJson(w, 200, logs)
		}
	})
}

func setCors(router *chi.Mux) {
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	connection, err := db.Connect()

	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	setCors(router)

	setDefaultRoutes(router, connection)

	// SETT MUTEX FOR LOG ERROR HANDLER
	var ErrorMu sync.Mutex
	adapters.ErrorMu = &ErrorMu

	router.Mount("/exchange", routes.ExchangeRoutes(connection))

	port := os.Getenv("PORT")
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server listing on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		panic("Error while trying to start the server")
	}
}
