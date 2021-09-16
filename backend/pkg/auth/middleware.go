package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const UserContextKey = contextKey("User")
const ColorContextKey = contextKey("Color")

type AnonUser struct {
	Id       string `json:"id"`
	Name     string `json:"user"`
	Color    string `json:"color"`
	Password string `json:"password"`
}

func (user *AnonUser) GetId() string {
	return user.Id
}

func (user *AnonUser) GetName() string {
	return user.Name
}
func (user *AnonUser) GetColor() string {
	return user.Color
}

func (user *AnonUser) GetPassword() string {
	return user.Password
}
func AuthMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, tok := r.URL.Query()["bearer"]
		name, nok := r.URL.Query()["user"]
		log.Println("token and name")
		log.Println(token, name)

		if tok && len(token) == 1 {
			user, err := ValidateToken(token[0])
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				log.Println("user and ctx")
				log.Println(user)
				log.Println(ctx)
				f(w, r.WithContext(ctx))
			}
		} else if nok && len(name) == 1 {
			user := AnonUser{Id: uuid.New().String(), Name: name[0]}
			ctx := context.WithValue(r.Context(), UserContextKey, &user)
			log.Println("user and ctx2")
			log.Println(user)
			log.Println(ctx)
			f(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Please login or provide name"))
		}
	})
}
