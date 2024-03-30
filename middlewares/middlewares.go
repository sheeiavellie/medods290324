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

		usersCollection := db.Database("amogus").Collection("users")

		var user *data.User
		err := usersCollection.FindOne(
			context.TODO(),
			bson.M{"user_id": userIDStr},
		).Decode(&user)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				log.Println("no document :(")
			default:
				log.Printf("error getting user: %v.", err)
			}
			return
		}

		// check if user exists in db
		// if he does, than issue tokens
		// then create refresh session
		next(w, r)
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
