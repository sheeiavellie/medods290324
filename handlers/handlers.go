package handlers

import (
	"log"
	"net/http"

	"github.com/sheeiavellie/medods290324/middlewares"
	"github.com/sheeiavellie/medods290324/util"
)

func HandleSingIn(w http.ResponseWriter, r *http.Request) {
	tokens, err := middlewares.GetTokensKey(r.Context())
	if err != nil {
		log.Printf("Error while getting tokens: %v.", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, tokens)
	if err != nil {
		log.Printf("Error while writing: %v.", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	tokens, err := middlewares.GetTokensKey(r.Context())
	if err != nil {
		log.Printf("Error while getting tokens: %v.", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, tokens)
	if err != nil {
		log.Printf("Error while writing: %v.", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
