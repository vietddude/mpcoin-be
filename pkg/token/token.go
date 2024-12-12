package token

import (
	"context"
	"errors"
	"fmt"
	"mpc/internal/db/redis"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

var (
	secretKey            = []byte("mpc")
	accessTokenDuration  = 12 * time.Hour
	refreshTokenDuration = 7 * 24 * time.Hour
)

type TokenManager struct {
	redis *redis.Client
}

func NewTokenManager(redis *redis.Client) *TokenManager {
	return &TokenManager{redis: redis}
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT for the specified token type and user ID.
func (tm *TokenManager) GenerateToken(ctx context.Context, userID uuid.UUID, tokenType string) (string, error) {
	// Revoke existing token first
	if err := tm.RevokeToken(ctx, userID, tokenType); err != nil {
		return "", fmt.Errorf("failed to revoke existing token: %w", err)
	}

	var duration time.Duration

	switch tokenType {
	case TokenTypeAccess:
		duration = accessTokenDuration
	case TokenTypeRefresh:
		duration = refreshTokenDuration
	default:
		return "", errors.New("invalid token type")
	}

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Store new token in Redis
	key := fmt.Sprintf("token:%s:%s", tokenType, userID.String())
	if err := tm.redis.Set(ctx, key, tokenString, duration).Err(); err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %w", err)
	}

	return tokenString, nil
}

// VerifyToken validates a JWT and extracts the user ID if valid.
func (tm *TokenManager) VerifyToken(ctx context.Context, tokenString, tokenType string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	// Check if token exists in Redis
	key := fmt.Sprintf("token:%s:%s", tokenType, claims.UserID.String())
	storedToken, err := tm.redis.Get(ctx, key).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to check token in Redis: %w", err)
	}
	if storedToken != tokenString {
		return uuid.Nil, errors.New("token has been revoked")
	}

	return claims.UserID, nil
}

// GenerateTokenPair creates both access and refresh tokens for a user ID.
func (tm *TokenManager) GenerateTokenPair(ctx context.Context, userID uuid.UUID) (Token, error) {
	accessToken, err := tm.GenerateToken(ctx, userID, TokenTypeAccess)
	if err != nil {
		return Token{}, err
	}

	refreshToken, err := tm.GenerateToken(ctx, userID, TokenTypeRefresh)
	if err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RevokeToken invalidates a token by removing it from Redis
func (tm *TokenManager) RevokeToken(ctx context.Context, userID uuid.UUID, tokenType string) error {
	key := fmt.Sprintf("token:%s:%s", tokenType, userID.String())
	if err := tm.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// RevokeAllUserTokens invalidates all tokens for a user
func (tm *TokenManager) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("token:%s:*", userID.String())
	keys, err := tm.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	if len(keys) > 0 {
		if err := tm.redis.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to revoke user tokens: %w", err)
		}
	}
	return nil
}
