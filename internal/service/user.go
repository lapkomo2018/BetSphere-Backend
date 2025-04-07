package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"stavki/external/hash"
	"stavki/internal/database"
	"stavki/internal/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type User struct {
	userDB *database.User
	rdb    *redis.Client
	hasher hash.Hasher
	auth   *Auth
}

// NewUser creates a new User service instance.
// UserDB must not have an active transaction.
func NewUser(userDB *database.User, rdb *redis.Client, h hash.Hasher, a *Auth) (*User, error) {
	if userDB.HasTx() {
		return nil, errors.New("userDB has transaction")
	}

	return &User{
		userDB: userDB,
		rdb:    rdb,
		hasher: h,
		auth:   a,
	}, nil
}

// Register registers a new user in the database.
func (u *User) Register(ctx context.Context, username, email, password string) (*model.User, model.TokenPair, error) {
	user := &model.User{
		Username: username,
		Email:    email,
		Password: u.hasher.Hash(password),
	}

	user, err := u.userDB.Create(ctx, user)
	if err != nil {
		return nil, model.TokenPair{}, err
	}

	u.cacheUser(ctx, user)

	pair, err := u.auth.CreateJWT(ctx, user.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    user.ID,
		}).Error("Error creating token pair")
		return nil, model.TokenPair{}, err
	}

	return user, pair, nil
}

// Login logs in a user and returns the user and token pair.
func (u *User) Login(ctx context.Context, login, password string) (*model.User, model.TokenPair, error) {
	user, err := u.userDB.GetByLogin(ctx, login)
	if err != nil {
		return nil, model.TokenPair{}, err
	}

	if !u.hasher.Compare(user.Password, password) {
		return nil, model.TokenPair{}, errors.New("invalid credentials")
	}

	pair, err := u.auth.CreateJWT(ctx, user.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    user.ID,
		}).Error("Error creating token pair")
		return nil, model.TokenPair{}, err
	}

	return user, pair, nil
}

func (u *User) Get(ctx context.Context, id uint64) (*model.User, error) {
	user, err := u.getUserFromCache(ctx, id)
	if err == nil {
		return user, nil
	} else if !errors.Is(err, redis.Nil) {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Error getting user from cache")
	}

	user, err = u.userDB.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	u.cacheUser(ctx, user)
	return user, nil
}

func (u *User) cacheUser(ctx context.Context, user *model.User) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    user.ID,
		}).Error("Error marshalling user")
		return
	}

	key := fmt.Sprintf("user:%d", user.ID)
	if err := u.rdb.Set(ctx, key, userJSON, time.Minute).Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    user.ID,
		}).Error("Error setting user in cache")
	}
}

func (u *User) getUserFromCache(ctx context.Context, id uint64) (*model.User, error) {
	key := fmt.Sprintf("user:%d", id)
	userCache, err := u.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(userCache), &user); err != nil {
		return nil, err
	}

	return &user, nil
}
