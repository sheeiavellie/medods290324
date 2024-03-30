package data

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID     uuid.UUID
	IsCool bool
}

type RefreshSession struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	ExpiresIn    time.Time
	CreatedAt    time.Time
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
