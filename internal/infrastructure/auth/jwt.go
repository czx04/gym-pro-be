package auth

import (
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	userdomain "gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    uuid.UUID       `json:"user_id"`
	Email     string          `json:"email"`
	Role      userdomain.Role `json:"role"`
	TokenType TokenType       `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret             string
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:             cfg.Secret,
		accessTokenExpire:  cfg.AccessTokenExpire,
		refreshTokenExpire: cfg.RefreshTokenExpire,
	}
}

func (j *JWTManager) GenerateAccessToken(userID uuid.UUID, email string, role userdomain.Role) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessTokenExpire)

	claims := Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

func (j *JWTManager) GenerateRefreshToken(userID uuid.UUID, email string, role userdomain.Role) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenExpire)

	claims := Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

func (j *JWTManager) GenerateTokenPair(userID uuid.UUID, email string, role userdomain.Role) (accessToken, refreshToken string, err error) {
	accessToken, err = j.GenerateAccessToken(userID, email, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.GenerateRefreshToken(userID, email, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, errors.TokenInvalid()
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.TokenInvalid()
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.TokenExpired()
	}

	return claims, nil
}

func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessToken {
		return nil, errors.TokenInvalid()
	}

	return claims, nil
}

func (j *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.TokenInvalid()
	}

	return claims, nil
}
