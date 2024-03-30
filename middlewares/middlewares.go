package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/sheeiavellie/medods290324/data"
	"github.com/sheeiavellie/medods290324/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SingIn(
	ctx context.Context,
	db *mongo.Client,
	tokenService *services.TokenService,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("userID")
		if userIDStr == "" {
			log.Printf("userID was not provided.")
			http.Error(w, "userID was not provided", http.StatusBadRequest)
			return
		}

		var user *data.User

		usersCollection := db.Database("amogus").Collection("users")
		err := usersCollection.FindOne(
			context.TODO(),
			bson.M{"user_id": userIDStr},
		).Decode(&user)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				log.Printf("User with ID: %s doen't exist: %v.", userIDStr, err)
				http.Error(w, "User with given ID doen't exist.", http.StatusBadRequest)
			default:
				log.Printf("Error getting user: %v.", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		tokens, err := tokenService.IssueTokens(user)
		if err != nil {
			log.Printf("Error while issueing tokens: %v.", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		session, err := data.NewRefreshSession(user, tokens.RefreshToken)
		if err != nil {
			log.Printf("Error while creating sessions: %v.", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sessionsCollection := db.Database("amogus").Collection("sessions")
		_, err = sessionsCollection.InsertOne(
			context.TODO(),
			session,
		)
		if err != nil {
			log.Printf("Error while creating session: %v.", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ctxRequestValue := context.WithValue(
			r.Context(),
			"",
			tokens,
		)
		next(w, r.WithContext(ctxRequestValue))
	}
}

func Refresh(
	ctx context.Context,
	db *mongo.Client,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
