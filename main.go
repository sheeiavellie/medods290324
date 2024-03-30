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
	port := "1337"

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
	log.Printf("server is listening on port: %s", port)
	http.ListenAndServe(":"+port, mux)

}
