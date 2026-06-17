package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	UserID        string
	WalletAddress string
	Role          string
	ProjectID     string
	IsOnboarded   bool
}

type TokenPair struct {
	RefreshToken string
	AccessToken  string
	TokenID      string
}

type AccessClaims struct {
	UserID        string `json:"uid"`
	WalletAddress string `json:"wallet"`
	Role          string `json:"role"`
	ProjectID     string `json:"pid,omitempty"`
	IsOnboarded   bool   `json:"onboarded"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID  string `json:"uid"`
	TokenID string `json:"tid"`
	jwt.RegisteredClaims
}

type Config struct {
	AccessTokenSecret   string
	RefreshTokenSecret  string
	AccessExpiryMinutes int
	RefreshExpiryHours  int
}

//generate access and refresh tokens

func GenerateTokenPair(payload TokenPayload, cfg Config) (*TokenPair, error) {

	accessExpiry := time.Duration(cfg.AccessExpiryMinutes) * time.Minute
	accessClaims := &AccessClaims{
		UserID:        payload.UserID,
		WalletAddress: payload.WalletAddress,
		Role:          payload.Role,
		ProjectID:     payload.ProjectID,
		IsOnboarded:   payload.IsOnboarded,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.AccessTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("error signing access token: %w", err)
	}

	tokenID := uuid.NewString()

	refreshExpiry := time.Duration(cfg.RefreshExpiryHours) * time.Hour
	refreshClaims := &RefreshClaims{
		UserID:  payload.UserID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(cfg.RefreshTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("signing refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenID:      tokenID,
	}, nil
}

//parse access Token

func ParseAccessToken(tokenString string, accessKey []byte) (*AccessClaims, error) {

	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return accessKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

//parse refresh token

func ParseRefreshToken(tokenString string, refreshKey []byte) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return refreshKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

//set cookies

func SetTokenCookies(c fiber.Ctx, pair *TokenPair, accessExpiryMinutes int, refreshExpiryHours int, isProd bool) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    pair.AccessToken,
		HTTPOnly: true,
		Secure:   isProd,
		SameSite: "Lax",
		MaxAge:   accessExpiryMinutes * 60,
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		HTTPOnly: true,
		Secure:   isProd,
		SameSite: "Lax",
		MaxAge:   refreshExpiryHours * 3600,
		Path:     "/",
	})
}

//clear cookies


func ClearTokenCookies(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{Name: "access_token", Value: "", HTTPOnly: true, MaxAge: -1, Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: "", HTTPOnly: true, MaxAge: -1, Path: "/"})
}