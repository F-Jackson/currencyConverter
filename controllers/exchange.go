package controllers

import (
	"errors"
	"fmt"
	"genesisbankly/exchange/db"
	"genesisbankly/exchange/utils"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	USD = "$"
	BRL = "R$"
	EUR = "€"
	BTC = "₿"
)

var ALLOWEDEXCHANGES = []string{"USD-BRL", "BRL-USD", "BRL-EUR", "EUR-BRL", "BTC-BRL", "BTC-USD"}

func getSymbol(currency string) string {
	switch currency {
	case "USD":
		return USD
	case "BRL":
		return BRL
	case "EUR":
		return EUR
	case "BTC":
		return BTC
	default:
		return ""
	}
}

func ConvertExchange(amount, rate float64, from, to string) (valueConverted float64, symbol string, err error) {
	valueConverted = amount * rate
	toCurrency := strings.ToUpper(to)
	fromCurrency := strings.ToUpper(from)

	exchange := fmt.Sprintf("%s-%s", toCurrency, fromCurrency)

	if !utils.Contains(ALLOWEDEXCHANGES, exchange) {
		return 0, "", errors.New("invalid currency codes exchange")
	}

	symbol = getSymbol(toCurrency)

	return valueConverted, symbol, nil
}

type ExchangeView struct {
	Id           uint64    `json:"id"`
	ViewedAt     time.Time `json:"viewed_at"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	FromValue    float64   `json:"from_value"`
	ToValue      float64   `json:"to_value"`
	TotalValue   float64   `json:"total_value"`
}

func handleGetExchangeError(err error) ([]ExchangeView, error) {
	return []ExchangeView{}, err
}

func ListExchanges(connection *gorm.DB) ([]ExchangeView, error) {
	exchanges, err := db.GetExchangesInsideHistory(connection)

	var newExchanges []ExchangeView

	if err != nil {
		return newExchanges, err
	}

	for _, exchange := range exchanges {
		fromValueStr, err := utils.DecryptGCM(string(exchange.FromValue))
		if err != nil {
			return handleGetExchangeError(err)
		}

		toValueStr, err := utils.DecryptGCM(string(exchange.ToValue))
		if err != nil {
			return handleGetExchangeError(err)
		}

		fromCurrency, err := utils.DecryptGCM(string(exchange.FromCurrency))
		if err != nil {
			return handleGetExchangeError(err)
		}

		toCurrency, err := utils.DecryptGCM(string(exchange.ToCurrency))
		if err != nil {
			return handleGetExchangeError(err)
		}

		totalValueStr, err := utils.DecryptGCM(string(exchange.TotalValue))
		if err != nil {
			return handleGetExchangeError(err)
		}

		totalValue, err := strconv.ParseFloat(totalValueStr, 64)
		if err != nil {
			return handleGetExchangeError(err)
		}

		toValue, err := strconv.ParseFloat(toValueStr, 64)
		if err != nil {
			return handleGetExchangeError(err)
		}

		fromValue, err := strconv.ParseFloat(fromValueStr, 64)
		if err != nil {
			return handleGetExchangeError(err)
		}

		newExchange := ExchangeView{
			Id: exchange.ID, ViewedAt: exchange.ViewedAt,
			ToCurrency: toCurrency, TotalValue: totalValue,
			FromCurrency: fromCurrency, ToValue: toValue,
			FromValue: fromValue,
		}

		newExchanges = append(newExchanges, newExchange)
	}

	return newExchanges, nil
}
