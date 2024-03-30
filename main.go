package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sheeiavellie/medods290324/handlers"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func main() {
	port := "1337"

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /sing-in/{userId}",
		handlers.HandleSingIn,
	)

	mux.HandleFunc(
		"POST /refresh",
		handlers.HandleRefresh,
	)

	log.Printf("server is listening on port: %s", port)
	http.ListenAndServe(":"+port, mux)

}
