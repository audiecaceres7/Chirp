package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"chirpy.com/database"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token_key := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)

	type new_user_params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := jwt.ParseWithClaims(token_key, jwt.MapClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
        fmt.Printf("Error parsing token: %s\n", err.Error())
		return
	}

    if !token.Valid {
        RespondWithError(w, http.StatusUnauthorized, "Time sesstion expired")
        return
    }

	id, err := token.Claims.GetSubject()
	if err != nil {
		fmt.Println(err.Error())
	}

    decoder := json.NewDecoder(r.Body)
    params := database.User{}
    decoder.Decode(&params)
  
    int_id, _ := strconv.Atoi(id)
	cfg.db.UpdateUser(int_id, params)

    RespondWithJSON(w, http.StatusOK, new_user_params{
        Email: params.Email, 
        Password: params.Password,
    }) 
}
