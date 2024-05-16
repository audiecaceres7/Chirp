package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password          string `json:"password"`
		Email             string `json:"email"`
		Expire_in_seconds int    `json:"expires_in_seconds"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Couldn't decode json: %v", err)
		return
	}

	all_users, _ := cfg.db.GetUsers()
	for _, user := range all_users {
		if user.Email == params.Email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Incorrect password. Please try again")
			} else {
                if params.Expire_in_seconds <= 0 {
                    params.Expire_in_seconds = 86000
                }

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
					Issuer:    "chirpy",
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(params.Expire_in_seconds))),
					Subject:   fmt.Sprint(user.ID),
				})

				t, _ := token.SignedString([]byte(cfg.jwtSecret))

                w.Header().Add("Authorization", fmt.Sprintf("Bearer %v", token)) 
				RespondWithJSON(w, http.StatusOK, UserLoginResponse{
					ID:    user.ID,
					Email: user.Email,
					Token: t,
				})
                return
			}
		}
	}
    RespondWithError(w, http.StatusUnauthorized, "Account doesn't exist")
}
