package routes

import (
	"genesisbankly/exchange/connectors"
	"genesisbankly/exchange/controllers"
	"genesisbankly/exchange/routes/adapters"
	"genesisbankly/exchange/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ReturnableExchange struct {
	ValorConvertido float64 `json:"valorConvertido"`
	SimboloMoeda    string  `json:"simboloMoeda"`
}

func handleMakeError(w http.ResponseWriter, r *http.Request, connection *gorm.DB, err error) {
	adapters.RespondWithError(w, connection, 402, err, r)
}

func handleMakeExchange(w http.ResponseWriter, r *http.Request, connection *gorm.DB) {
	amount := chi.URLParam(r, "amount")
	from := chi.URLParam(r, "from")
	to := chi.URLParam(r, "to")
	rate := chi.URLParam(r, "rate")

	err := utils.VerifyUrlParamsNullabity(amount, from, to, rate)
	if err != nil {
		handleMakeError(w, r, connection, err)
		return
	}

	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		handleMakeError(w, r, connection, err)
		return
	}

	floatRate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		handleMakeError(w, r, connection, err)
		return
	}

	totalValue, symbol, err := connectors.ConvertAndSaveExchange(floatAmount, floatRate, from, to, connection)
	if err != nil {
		handleMakeError(w, r, connection, err)
		return
	}

	adapters.RespondWithJson(w, 500, &ReturnableExchange{
		ValorConvertido: totalValue,
		SimboloMoeda:    symbol,
	})
}

func handleListExchanges(w http.ResponseWriter, r *http.Request, connection *gorm.DB) {
	exchanges, err := controllers.ListExchanges(connection)

	if err != nil {
		adapters.RespondWithError(w, connection, 500, err, r)
	} else {
		adapters.RespondWithJson(w, 200, exchanges)
	}
}

func ExchangeRoutes(connection *gorm.DB) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/{amount}/{from}/{to}/{rate}", func(w http.ResponseWriter, r *http.Request) {
		handleMakeExchange(w, r, connection)
	})

	router.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		handleListExchanges(w, r, connection)
	})
	return router
}
