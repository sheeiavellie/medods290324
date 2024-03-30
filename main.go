package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sheeiavellie/medods290324/handlers"
	"github.com/sheeiavellie/medods290324/middlewares"
)

func main() {
	port := "1337"
	ctx := context.Background()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /sing-in",
		middlewares.ValidateUser(
			ctx,
			middlewares.SingIn(
				ctx,
				handlers.HandleSingIn,
			),
		),
	)

	mux.HandleFunc(
		"POST /refresh",
		middlewares.ValidateRefresh(
			ctx,
			middlewares.Refresh(
				ctx,
				handlers.HandleRefresh,
			),
		),
	)

	log.Printf("server is listening on port: %s", port)
	http.ListenAndServe(":"+port, mux)

}
