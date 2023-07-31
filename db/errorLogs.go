package db

import (
	"time"

	"gorm.io/gorm"
)

type ErrorLog struct {
	ID       uint64    `gorm:"->;primaryKey;auto_increment"`
	Message  string    `gorm:"not null"`
	Code     int       `gorm:"not null"`
	LoggedAt time.Time `gorm:"not null"`
	Payload  []byte    `gorm:"not null"`
}

func ConnectErrorLog(connection *gorm.DB) error {
	err := connection.AutoMigrate(&ErrorLog{})
	return err
}

func CreateLogInsideDb(connection *gorm.DB, msg string, code int, payload []byte) (ErrorLog, error) {
	log := ErrorLog{
		Message: msg, Code: code,
		Payload: payload, LoggedAt: time.Now().UTC(),
	}

	result := connection.Create(&log)

	if result.Error != nil {
		return ErrorLog{}, result.Error
	}
	return log, nil
}

func ListLogsInsideDb(connection *gorm.DB) ([]ErrorLog, error) {
	logs := []ErrorLog{}
	result := connection.Find(&logs)

	if result.Error != nil {
		return []ErrorLog{}, result.Error
	}
	return logs, nil
}
