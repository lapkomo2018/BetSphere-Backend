package database

import (
	"errors"

	"gorm.io/gorm"
)

type (
	TransactionProvider interface {
		Transact(txFunc func(Adapters) error) error
	}
	transactionProvider struct {
		db *gorm.DB
	}
	Adapters struct {
		UserRepository *UserRepository
		JWTRepository  *JWTRepository
	}
)

func NewTransactionProvider(db *gorm.DB) TransactionProvider {
	return &transactionProvider{
		db: db,
	}
}

func (p *transactionProvider) Transact(txFunc func(adapters Adapters) error) error {
	return runInTx(p.db, func(tx *gorm.DB) error {
		adapters := Adapters{
			UserRepository: NewUserRepository(tx),
			JWTRepository:  NewJWTRepository(tx),
		}

		return txFunc(adapters)
	})
}

func runInTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := fn(tx)
	if err == nil {
		return tx.Commit().Error
	}

	rollbackErr := tx.Rollback().Error
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
