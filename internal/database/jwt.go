package database

import (
	"context"

	"stavki/internal/model"

	"gorm.io/gorm"
)

type JWTRepository struct {
	db *gorm.DB
}

func NewJWTRepository(db *gorm.DB) *JWTRepository {
	return &JWTRepository{
		db: db,
	}
}

func (j *JWTRepository) Create(ctx context.Context, jwt *model.JWT) (*model.JWT, error) {
	return jwt, j.db.WithContext(ctx).Create(jwt).Error
}

func (j *JWTRepository) Get(ctx context.Context, token string) (*model.JWT, error) {
	var jwt model.JWT
	return &jwt, j.db.WithContext(ctx).Where("refresh_token = ?", token).First(&jwt).Error
}

func (j *JWTRepository) Delete(ctx context.Context, token string) error {
	return j.db.WithContext(ctx).Where("refresh_token = ?", token).Delete(&model.JWT{}).Error
}

func (j *JWTRepository) DeleteExpired(ctx context.Context) (count int64, err error) {
	res := j.db.WithContext(ctx).Where("expires_at < ?", gorm.Expr("NOW()")).Delete(&model.JWT{})
	return res.RowsAffected, res.Error
}
