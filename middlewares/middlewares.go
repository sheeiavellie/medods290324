package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sheeiavellie/medods290324/data"
	"github.com/sheeiavellie/medods290324/services"
	"github.com/sheeiavellie/medods290324/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func SingIn(
	ctx context.Context,
	sessionService *services.SessionService,
	tokenService *services.TokenService,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userID")
		if userID == "" {
			log.Printf("userID was not provided.")
			http.Error(w, "userID was not provided.", http.StatusBadRequest)
			return
		}

		user, err := sessionService.GetUser(context.TODO(), userID)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				log.Printf("User with ID: %s doen't exist: %v.", userID, err)
				http.Error(w, "User with given ID doen't exist.", http.StatusBadRequest)
			default:
				log.Printf("Error getting user: %v.", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		refreshSessionID := uuid.New().String()

		tokens, err := tokenService.IssueTokens(user.ID, refreshSessionID)
		if err != nil {
			log.Printf("Error while issueing tokens: %v.", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = sessionService.CreateRefreshSession(
			context.TODO(),
			refreshSessionID,
			user.ID,
			tokens.RefreshToken,
		)
		if err != nil {
			log.Printf("Error while creating sessions: %v.", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ctxRequestValue := context.WithValue(
			r.Context(),
			contextKeyTokens,
			tokens,
		)
		next(w, r.WithContext(ctxRequestValue))
	}
}

func Refresh(
	ctx context.Context,
	sessionService *services.SessionService,
	tokenService *services.TokenService,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshRequest data.RefreshRequest
		err := util.ReadJSON(r, &refreshRequest)
		if err != nil {
			log.Printf("Refresh token wasn't provided: %v.", err)
			http.Error(w, "Refresh token wasn't provided.", http.StatusBadRequest)
			return
		}

		refreshTokenClaims, err := tokenService.DecodeRefreshToken(
			refreshRequest.RefreshToken,
		)
		if err != nil {
		}

		err = sessionService.ValidateRefreshSession(
			context.TODO(),
			refreshTokenClaims.SessionID,
			refreshRequest.RefreshToken,
		)
		if err != nil {
		}

		err = sessionService.DeleteRefreshSession(
			context.TODO(),
			refreshTokenClaims.SessionID,
		)
		if err != nil {
		}

		refreshSessionID := uuid.New().String()

		tokens, err := tokenService.IssueTokens(
			refreshTokenClaims.UserID,
			refreshSessionID,
		)
		if err != nil {
			log.Printf("Error while issueing tokens: %v.", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = sessionService.CreateRefreshSession(
			context.TODO(),
			refreshSessionID,
			refreshTokenClaims.UserID,
			tokens.RefreshToken,
		)
		if err != nil {
			log.Printf("Error while creating sessions: %v.", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ctxRequestValue := context.WithValue(
			r.Context(),
			contextKeyTokens,
			tokens,
		)
		next(w, r.WithContext(ctxRequestValue))
	}
}
