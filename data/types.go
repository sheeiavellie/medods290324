package data

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTTL = time.Hour * 24 * 30
	cost       = 10
)

type User struct {
	ID     uuid.UUID `bson:"user_id"`
	IsCool bool      `bson:"is_cool"`
}

func (u *User) UnmarshalBSON(data []byte) error {
	var raw bson.M
	if err := bson.Unmarshal(data, &raw); err != nil {
		return err
	}

	idStr, ok := raw["user_id"].(string)
	if !ok {
		return fmt.Errorf("user_id field not found or not a string")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	isCool, ok := raw["is_cool"].(bool)
	if !ok {
		return fmt.Errorf("is_cool field not found or not a boolean")
	}

	u.ID = id
	u.IsCool = isCool

	return nil
}

type RefreshSession struct {
	UserID       uuid.UUID `bson:"user_id"`
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
