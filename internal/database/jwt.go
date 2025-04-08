package database

import (
	"context"

	"stavki/internal/model"
)

type JWT struct {
	*Base[JWT]
}

func NewJWT(db Database, tx *Transaction) *JWT {
	return &JWT{
		Base: NewBase[JWT](db, tx, NewJWT),
	}
}

func (j *JWT) Create(ctx context.Context, jwt *model.JWT) (*model.JWT, error) {
	tx := j.tx
	if tx == nil {
		tx = j.db.TransactionWithContext(ctx)
		defer tx.EnsureRollback()
	}

	if err := tx.Create(jwt).Error; err != nil {
		return jwt, err
	}

	if j.tx == nil {
		return jwt, tx.Commit().Error
	}

	return jwt, nil
}

func (j *JWT) Get(ctx context.Context, token string) (*model.JWT, error) {
	tx := j.tx
	if tx == nil {
		tx = j.db.TransactionWithContext(ctx)
		defer tx.EnsureRollback()
	}

	var jwt model.JWT
	if err := tx.Where("refresh_token = ?", token).First(&jwt).Error; err != nil {
		return &jwt, err
	}

	if j.tx == nil {
		return &jwt, tx.Commit().Error
	}

	return &jwt, nil
}

func (j *JWT) Delete(ctx context.Context, token string) error {
	tx := j.tx
	if tx == nil {
		tx = j.db.TransactionWithContext(ctx)
		defer tx.EnsureRollback()
	}

	if err := tx.Where("refresh_token = ?", token).Delete(&model.JWT{}).Error; err != nil {
		return err
	}

	if j.tx == nil {
		return tx.Commit().Error
	}

	return nil
}
