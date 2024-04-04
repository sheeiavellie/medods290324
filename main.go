package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sheeiavellie/medods290324/handlers"
	"github.com/sheeiavellie/medods290324/middlewares"
	"github.com/sheeiavellie/medods290324/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Не использовал контексты с дэдлаином нигде, так как
// база стоит локально
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	port := os.Getenv("PORT")
	secret := os.Getenv("SECRET")
	connStr := os.Getenv("CONN_STRING")

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	tokenService := services.NewTokenService([]byte(secret))
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
