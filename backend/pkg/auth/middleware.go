package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const UserContextKey = contextKey("User")
const ColorContextKey = contextKey("Color")

type AnonUser struct {
	Id    string `json:"id"`
	User  string `json:"user"`
	Color string `json:"color"`
}

func (user *AnonUser) GetId() string {
	return user.Id
}

func (user *AnonUser) GetUser() string {
	return user.User
}
func (user *AnonUser) GetColor() string {
	return user.Color
}

func AuthMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, tok := r.URL.Query()["bearer"]
		user, nok := r.URL.Query()["user"]

		if tok && len(token) == 1 {
			user, err := ValidateToken(token[0])
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				f(w, r.WithContext(ctx))
			}
		} else if nok && len(user) == 1 {
			user := AnonUser{Id: uuid.New().String(), User: user[0]}
			ctx := context.WithValue(r.Context(), UserContextKey, &user)
			f(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Please login or provide name"))
		}
	})
}
