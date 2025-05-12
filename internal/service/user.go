package service

import (
	"context"
	"errors"

	"stavki/external/hash"
	"stavki/internal/cache"
	"stavki/internal/database"
	"stavki/internal/log"
	"stavki/internal/model"
)

type UserService struct {
	txProvider database.TransactionProvider
	userDB     *database.UserRepository
	r          *cache.Cache
	hasher     hash.Hasher
	auth       *Auth
}

// NewUser creates a new UserService service instance.
// UserDB must not have an active transaction.
func NewUser(txProvider database.TransactionProvider, userDB *database.UserRepository, r *cache.Cache, h hash.Hasher, a *Auth) *UserService {
	return &UserService{
		txProvider: txProvider,
		userDB:     userDB,
		r:          r,
		hasher:     h,
		auth:       a,
	}
}

// Register registers a new user in the database.
func (u *UserService) Register(ctx context.Context, username, email, password string) (*model.User, model.TokenPair, error) {
	user := &model.User{
		Username: username,
		Email:    email,
		Password: u.hasher.Hash(password),
	}

	user, err := u.userDB.Create(ctx, user)
	if err != nil {
		return nil, model.TokenPair{}, err
	}

	if err := u.r.SetUser(ctx, user); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    user.ID,
		}).Error("Error caching user")
	}

	pair, err := u.auth.CreateJWT(ctx, user.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    user.ID,
		}).Error("Error creating token pair")
		return nil, model.TokenPair{}, err
	}

	return user, pair, nil
}

// Login logs in a user and returns the user and token pair.
func (u *UserService) Login(ctx context.Context, login, password string) (*model.User, model.TokenPair, error) {
	user, err := u.userDB.GetByLogin(ctx, login)
	if err != nil {
		return nil, model.TokenPair{}, err
	}

	if !u.hasher.Compare(user.Password, password) {
		return nil, model.TokenPair{}, errors.New("invalid credentials")
	}

	pair, err := u.auth.CreateJWT(ctx, user.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    user.ID,
		}).Error("Error creating token pair")
		return nil, model.TokenPair{}, err
	}

	return user, pair, nil
}

func (u *UserService) Get(ctx context.Context, id uint64) (*model.User, error) {
	return cache.Run[*model.User](
		ctx,
		u.r.R(),
		cache.UserCacheKey(id),
		cache.UserCacheTTL,
		func(user *model.User) (*model.User, error) {
			if user != nil {
				return user, nil
			}

			return u.userDB.Get(ctx, id)
		})
}
