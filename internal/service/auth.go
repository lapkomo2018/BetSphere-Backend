package service

import (
	"context"
	"errors"
	"time"

	"stavki/internal/database"
	"stavki/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type (
	Auth struct {
		txProvider database.TransactionProvider
		jwtDB      *database.JWTRepository
		jwtKey     []byte
	}
)

var (
	// accessTokenTTL is the time-to-live for access tokens.
	accessTokenTTL = 15 * time.Minute
	// refreshTokenTTL is the time-to-live for refresh tokens.
	refreshTokenTTL = 30 * 24 * time.Hour

	// ErrInvalidToken is returned when the token is invalid.
	ErrInvalidToken = errors.New("invalid token")
)

// NewAuth creates a new Auth service instance.
// JWT_DB must not have an active transaction.
func NewAuth(txProvider database.TransactionProvider, jwtDB *database.JWTRepository, jwtKey []byte) (*Auth, error) {
	return &Auth{
		txProvider: txProvider,
		jwtDB:      jwtDB,
		jwtKey:     jwtKey,
	}, nil
}

// CreateJWT creates a new JWT for the given user ID.
// It generates a new access token and refresh token pair, saves it to the database,
// and returns the JWT object.
func (a *Auth) CreateJWT(ctx context.Context, userID uint64) (model.TokenPair, error) {
	pair, err := a.createTokenPair(userID)
	if err != nil {
		return model.TokenPair{}, err
	}

	if _, err := a.jwtDB.Create(ctx, &model.JWT{
		RefreshToken: pair.RefreshToken.Token,
		UserID:       userID,
		ExpiresAt:    pair.RefreshToken.ExpiresAt,
		CreatedAt:    time.Now(),
	}); err != nil {
		return model.TokenPair{}, err
	}

	return pair, nil
}

// AuthenticateJWT authenticates the given JWT.
// It checks if the access token is valid and not expired.
// If the access token is valid, it returns the user ID.
func (a *Auth) AuthenticateJWT(ctx context.Context, token string) (uint64, error) {
	// Parse the token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.jwtKey, nil
	})
	if err != nil {
		return 0, err
	}

	// Check if the token is valid and not expired
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userID := uint64(claims["sub"].(float64))
		return userID, nil
	}

	return 0, ErrInvalidToken
}

// RefreshJWT refreshes the given JWT.
func (a *Auth) RefreshJWT(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	// Parse the refresh token
	parsedToken, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.jwtKey, nil
	})
	if err != nil {
		return model.TokenPair{}, err
	}

	// Check if the token is valid and not expired
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return model.TokenPair{}, ErrInvalidToken
	}

	var newPair model.TokenPair
	return newPair, a.txProvider.Transact(func(adapters database.Adapters) error {
		// Check if the refresh token is in the database
		jwtSession, err := adapters.JWTRepository.Get(ctx, refreshToken)
		if err != nil {
			return err
		}

		// Check if the user ID in the token matches the user ID in the database
		if jwtSession.UserID != uint64(claims["sub"].(float64)) {
			return errors.New("user ID mismatch")
		}

		// Generate a new access token and refresh token pair
		newPair, err = a.createTokenPair(jwtSession.UserID)
		if err != nil {
			return err
		}

		// Update the refresh token in the database
		if err := adapters.JWTRepository.Delete(ctx, refreshToken); err != nil {
			return err
		}

		if _, err := adapters.JWTRepository.Create(ctx, &model.JWT{
			RefreshToken: newPair.RefreshToken.Token,
			UserID:       jwtSession.UserID,
			ExpiresAt:    newPair.RefreshToken.ExpiresAt,
			CreatedAt:    time.Now(),
		}); err != nil {
			return err
		}

		return nil
	})
}

// LogoutJWT logs out the user by deleting the refresh token from the database.
func (a *Auth) LogoutJWT(ctx context.Context, refreshToken string) error {
	if err := a.jwtDB.Delete(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

// createTokenPair generates a new access token and refresh token pair.
func (a *Auth) createTokenPair(userID uint64) (model.TokenPair, error) {
	// Generate access token and refresh token
	accessToken, err := a.generateJWTToken(a.jwtKey, userID, accessTokenTTL)
	if err != nil {
		return model.TokenPair{}, err
	}

	refreshToken, err := a.generateJWTToken(a.jwtKey, userID, refreshTokenTTL)
	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateJWTToken generates a new jwt token for the given user ID.
func (a *Auth) generateJWTToken(key []byte, userID uint64, ttl time.Duration) (model.JWTToken, error) {
	expiresAt := time.Now().Add(ttl)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expiresAt.Unix(),
	})

	token, err := t.SignedString(key)
	if err != nil {
		return model.JWTToken{}, err
	}

	return model.JWTToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
