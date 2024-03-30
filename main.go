package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /amogus",
		func(w http.ResponseWriter, _ *http.Request) {
			err := WriteJSON(w, http.StatusOK, "AMONGAS")
			if err != nil {
				log.Printf("error: %v.", err)
				return
			}
		},
	)

}
