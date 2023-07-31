package connectors

import (
	"fmt"
	"genesisbankly/exchange/controllers"
	"genesisbankly/exchange/db"
	"genesisbankly/exchange/utils"

	"gorm.io/gorm"
)

func handleConvertError(err error) (float64, string, error) {
	return 0, "", err
}

func ConvertAndSaveExchange(floatAmount, floatRate float64, from, to string, connection *gorm.DB) (float64, string, error) {
	totalValue, symbol, err := controllers.ConvertExchange(floatAmount, floatRate, from, to)
	if err != nil {
		return handleConvertError(err)
	}

	encryptedFloatAmount, err := utils.EncryptGCM(fmt.Sprintf("%f", floatAmount))
	if err != nil {
		return handleConvertError(err)
	}

	encryptedFloatRate, err := utils.EncryptGCM(fmt.Sprintf("%f", floatRate))
	if err != nil {
		return handleConvertError(err)
	}

	encryptedTotal, err := utils.EncryptGCM(fmt.Sprintf("%f", totalValue))
	if err != nil {
		return handleConvertError(err)
	}

	encryptedFrom, err := utils.EncryptGCM(from)
	if err != nil {
		return handleConvertError(err)
	}

	encryptedTo, err := utils.EncryptGCM(to)
	if err != nil {
		return handleConvertError(err)
	}

	_, err = db.SaveExchangeInsideHistory(
		encryptedFloatAmount, encryptedFloatRate,
		encryptedTotal, encryptedFrom,
		encryptedTo, connection,
	)

	if err != nil {
		return handleConvertError(err)
	}

	return totalValue, symbol, err
}
