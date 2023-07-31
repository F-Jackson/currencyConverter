package adapters

import (
	"fmt"
	"genesisbankly/exchange/db"
	"genesisbankly/exchange/utils"
	"net/http"
	"sync"

	"gorm.io/gorm"
)

var ErrorMu *sync.Mutex

type ErrResponse struct {
	Message string `json:"message"`
}

func createErrorLog(connection *gorm.DB, code int, err error, payload interface{}) {
	ErrorMu.Lock()
	defer ErrorMu.Unlock()

	payloadStr := fmt.Sprintf("{%v}", payload)
	payloadEncrypted, newErr := utils.EncryptGCM(payloadStr)
	if newErr == nil {
		db.CreateLogInsideDb(connection, err.Error(), code, payloadEncrypted)
	}
}

func RespondWithError(w http.ResponseWriter, connection *gorm.DB, code int, err error, payload interface{}) {
	go createErrorLog(connection, code, err, payload)
	RespondWithJson(w, code, ErrResponse{Message: err.Error()})
}
