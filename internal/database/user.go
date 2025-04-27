package database

import (
	"context"
	"errors"

	"stavki/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return user, u.db.WithContext(ctx).Create(user).Error
}

func (u *UserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	return &user, u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
}

func (u *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	return &user, u.db.WithContext(ctx).Where("username = ? OR email = ?", login, login).First(&user).Error
}

func (u *UserRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	return user, u.db.WithContext(ctx).Save(user).Error
}

func (u *UserRepository) UpdateBalance(ctx context.Context, id uint64, amount float64) (*model.User, error) {
	var user model.User
	result := u.db.WithContext(ctx).
		Model(&user).
		Where("id = ?", id).
		Where("balance + ? >= 0", amount).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		return &user, result.Error
	}

	if result.RowsAffected == 0 {
		return &user, errors.New("insufficient funds")
	}

	return &user, nil
}
