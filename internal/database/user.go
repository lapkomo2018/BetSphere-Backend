package database

import (
	"context"

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
