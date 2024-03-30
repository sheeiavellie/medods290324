package data

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTTL = time.Hour * 24 * 30
	cost       = 10
)

type User struct {
	ID     string `bson:"user_id"`
	IsCool bool   `bson:"is_cool"`
}

type RefreshSession struct {
	UserID       string    `bson:"user_id"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiresIn    time.Time `bson:"expires_in"`
	CreatedAt    time.Time `bson:"created_at"`
}

func NewRefreshSession(user *User, refreshToken string) (*RefreshSession, error) {
	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 10)
	if err != nil {
		return nil, err
	}
	return &RefreshSession{
		UserID:       user.ID,
		RefreshToken: string(refreshTokenHash),
		ExpiresIn:    time.Now().Add(refreshTTL),
		CreatedAt:    time.Now(),
	}, nil
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
