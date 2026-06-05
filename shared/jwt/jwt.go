package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	UserID      string
	Role        string
	WorkSpaceID string
	Version     int64
	IsOnboarded bool
}

type TokenPair struct {
	RefreshToken string
	AccessToken  string
	TokenID      string
}

type AccessClaims struct {
	UserID      string `json:"uid"`
	WorkspaceID string `json:"wid"`
	Role        string `json:"role"`
	Version     int64  `json:"ver"`
	IsOnboarded bool   `json:"onboarded"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID  string `json:"uid"`
	TokenID string `json:"tid"`
	jwt.RegisteredClaims
}

type Config struct{
	AccessTokenSecret   string
	RefreshTokenSecret  string
	AccessExpiryMinutes int
	RefreshExpiryHours  int
}

func GenerateTokenPair(payload TokenPayload,cfg Config )(*TokenPair,error){

	//access Token
	 accesExpiry:=time.Duration(cfg.AccessExpiryMinutes)*time.Minute
	 accessClaims:=&AccessClaims{
		UserID: payload.UserID,
		Role: payload.Role,
		Version: payload.Version,
		WorkspaceID: payload.WorkSpaceID,
		IsOnboarded: payload.IsOnboarded,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accesExpiry)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	 }
	 accesToken,err:=jwt.NewWithClaims(jwt.SigningMethodHS256,accessClaims).SignedString([]byte(cfg.AccessTokenSecret))
	 if err != nil {
		return nil, fmt.Errorf("error signing access token: %w", err)
	}

	tokenID:=uuid.NewString()

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
		AccessToken:  accesToken,
		RefreshToken: refreshToken,
		TokenID:      tokenID,
	}, nil
}