package data

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID     string `bson:"user_id"`
	IsCool bool   `bson:"is_cool"`
}

type RefreshTokenClaims struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}

type RefreshSession struct {
	ID           string    `bson:"session_id"`
	UserID       string    `bson:"user_id"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiresIn    time.Time `bson:"expires_in"`
	IssuedAt     time.Time `bson:"issued_at"`
}

func NewRefreshSession(
	sessionID string,
	userID string,
	refreshToken string,
	issuingTime time.Time,
	refreshTTL time.Duration,
) (*RefreshSession, error) {
	refreshTokenHash, err := bcrypt.GenerateFromPassword(
		[]byte(refreshToken),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	return &RefreshSession{
		ID:           sessionID,
		UserID:       userID,
		RefreshToken: string(refreshTokenHash),
		ExpiresIn:    issuingTime.Add(refreshTTL),
		IssuedAt:     issuingTime,
	}, nil
}

type JWTTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
