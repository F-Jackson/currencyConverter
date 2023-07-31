package db

import (
	"time"

	"gorm.io/gorm"
)

type ExchangeHistory struct {
	ID           uint64    `gorm:"->;primaryKey;auto_increment"`
	ViewedAt     time.Time `gorm:"not null"`
	FromCurrency []byte    `gorm:"not null"`
	ToCurrency   []byte    `gorm:"not null"`
	FromValue    []byte    `gorm:"not null"`
	ToValue      []byte    `gorm:"not null"`
	TotalValue   []byte    `gorm:"not null"`
}

func ConnectExchange(connection *gorm.DB) error {
	err := connection.AutoMigrate(&ExchangeHistory{})
	return err
}

func SaveExchangeInsideHistory(fromValue, toValue, totalValue, fromCurrency, toCurrency []byte, connection *gorm.DB) (ExchangeHistory, error) {
	exchange := ExchangeHistory{
		FromCurrency: fromCurrency, ToValue: toValue,
		TotalValue: totalValue, FromValue: fromValue,
		ToCurrency: toCurrency, ViewedAt: time.Now().UTC(),
	}

	result := connection.Create(&exchange)

	if result.Error != nil {
		return ExchangeHistory{}, result.Error
	}
	return exchange, nil
}

func GetExchangesInsideHistory(connection *gorm.DB) ([]ExchangeHistory, error) {
	exchanges := []ExchangeHistory{}
	result := connection.Find(&exchanges)

	if result.Error != nil {
		return []ExchangeHistory{}, result.Error
	}
	return exchanges, nil
}
