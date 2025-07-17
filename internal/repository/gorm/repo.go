package gorm

import (
	"errors"

	"gorm.io/gorm"
)

func handleGormError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}

	// Laporkan error ke Sentry
	// sentry.CaptureException(err)

	return errors.New("sql_error")
}
