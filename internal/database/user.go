package database

import (
	"context"

	"stavki/internal/model"
)

type User struct {
	*Base[User]
}

func NewUser(db Database, tx *Transaction) *User {
	return &User{
		Base: NewBase[User](db, tx, NewUser),
	}
}

func (u *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	tx := u.tx
	if tx == nil {
		tx = u.db.TransactionWithContext(ctx)
		defer tx.EnsureRollback()
	}

	if err := tx.Create(user).Error; err != nil {
		return user, err
	}

	if u.tx == nil {
		return user, tx.Commit().Error
	}

	return user, nil
}

func (u *User) Get(ctx context.Context, id uint64) (*model.User, error) {
	tx := u.tx
	if tx == nil {
		tx = u.db.TransactionWithContext(ctx)
		defer tx.EnsureRollback()
	}

	var user model.User
	if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
		return &user, err
	}

	if u.tx == nil {
		return &user, tx.Commit().Error
	}

	return &user, nil
}

func (u *User) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	tx := u.tx
	if tx == nil {
		tx = u.db.TransactionWithContext(ctx)
		defer tx.EnsureRollback()
	}

	var user model.User
	if err := tx.Where("username = ? OR email = ?", login, login).First(&user).Error; err != nil {
		return &user, err
	}

	if u.tx == nil {
		return &user, tx.Commit().Error
	}

	return &user, nil
}

func (u *User) Save(ctx context.Context, user *model.User) (*model.User, error) {
	tx := u.tx
	if tx == nil {
		tx = u.db.TransactionWithContext(ctx)
		defer tx.EnsureRollback()
	}

	if err := tx.Save(user).Error; err != nil {
		return user, err
	}

	if u.tx == nil {
		return user, tx.Commit().Error
	}

	return user, nil
}
