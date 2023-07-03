package main

import (
	"encoding/json"
	"net/http"

	"github.com/bilalbg/GoLang-Webchat/backend/pkg/auth"
	"github.com/bilalbg/GoLang-Webchat/backend/pkg/repository"
)

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Color    string `json:"color"`
}

type API struct {
	UserRepository *repository.UserRepository
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser := api.UserRepository.FindUserByUsername(user.Username)
	if dbUser == nil {
		returnErrorResponse(w)
		return
	}

	ok, err := auth.ComparePassword(user.Password, dbUser.Password)
	if !ok || err != nil {
		returnErrorResponse(w)
		return
	}

	token, err := auth.CreateJWTToken(dbUser)

	if err != nil {
		returnErrorResponse(w)
		return
	}

	w.Write([]byte(token))
}

func returnErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"error\""))
}
