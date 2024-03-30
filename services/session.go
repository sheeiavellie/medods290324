package services

import (
	"context"

	"github.com/sheeiavellie/medods290324/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	user *data.User,
	refreshToken string,
) error {
	session, err := data.NewRefreshSession(user, refreshToken)
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

func (s *SessionService) GetRefreshSession(refreshToken string) {
}
