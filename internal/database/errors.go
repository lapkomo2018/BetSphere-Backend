package database

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
	ErrTransactionNil = errors.New("transaction is nil")
)
