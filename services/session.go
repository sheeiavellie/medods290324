package services

import (
	"context"
	"time"

	"github.com/sheeiavellie/medods290324/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTTL = time.Hour * 24 * 30
)

type SessionService struct {
	client *mongo.Client
}

func NewSessionService(client *mongo.Client) *SessionService {
	return &SessionService{
		client: client,
	}
}

func (s *SessionService) GetUser(
	ctx context.Context,
	userID string,
) (*data.User, error) {
	var user *data.User

	usersCollection := s.client.Database("amogus").Collection("users")

	err := usersCollection.FindOne(
		ctx,
		bson.M{"user_id": userID},
	).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SessionService) CreateRefreshSession(
	ctx context.Context,
	sessionID string,
	userID string,
	refreshToken string,
) error {
	session, err := data.NewRefreshSession(
		sessionID,
		userID,
		refreshToken,
		time.Now(),
		refreshTTL,
	)
	if err != nil {
		return err
	}

	sessionsCollection := s.client.Database("amogus").Collection("sessions")
	_, err = sessionsCollection.InsertOne(
		ctx,
		session,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionService) DeleteRefreshSession(
	ctx context.Context,
	sessionID string,
) error {
	sessionsCollection := s.client.Database("amogus").Collection("sessions")

	_, err := sessionsCollection.DeleteOne(
		ctx,
		bson.M{"session_id": sessionID},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionService) ValidateRefreshSession(
	ctx context.Context,
	sessionID string,
	refreshToken string,
) error {
	sessionsCollection := s.client.Database("amogus").Collection("sessions")

	var curSession data.RefreshSession
	err := sessionsCollection.FindOne(
		ctx,
		bson.M{"session_id": sessionID},
	).Decode(&curSession)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(curSession.RefreshToken),
		[]byte(refreshToken),
	)
	if err != nil {
		return err
	}

	if !time.Now().Before(curSession.ExpiresIn) {
		return err
	}

	return nil
}
