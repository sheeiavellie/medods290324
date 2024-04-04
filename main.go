package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sheeiavellie/medods290324/handlers"
	"github.com/sheeiavellie/medods290324/middlewares"
	"github.com/sheeiavellie/medods290324/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	port := "1337"
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:admin@localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	tokenService := services.NewTokenService([]byte("secret"))
	sessionService := services.NewSessionService(client)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /sing-in",
		middlewares.SingIn(
			ctx,
			sessionService,
			tokenService,
			handlers.HandleSingIn,
		),
	)

	mux.HandleFunc(
		"POST /refresh",
		middlewares.Refresh(
			ctx,
			sessionService,
			tokenService,
			handlers.HandleRefresh,
		),
	)

	log.Printf("server is listening on port: %s", port)
	http.ListenAndServe(":"+port, mux)

}
