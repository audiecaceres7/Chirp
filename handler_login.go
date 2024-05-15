package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
        Email    string `json:"email"`
		Password string `json:"password"`
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
                fmt.Println(err)
                return
            } else {
                RespondWithJSON(w, http.StatusOK, User{
                    ID:    user.ID,
                    Email: user.Email,
                })
            }
        }
	}


}
