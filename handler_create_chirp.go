package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func ReplaceProfaneWords(message string) string {
	profaneWords := []string{"kerfuffle", "fornax", "sharbert"}
	arr := strings.Split(message, " ")
	for i := range arr {
		for j := range profaneWords {
			if strings.ToLower(arr[i]) == strings.ToLower(profaneWords[j]) {
				arr[i] = "****"
			}
		}
	}
	return strings.Join(arr, " ")
}

func (cfg *apiConfig) HandleChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
        Body string `json:"body"`
	}

    decoder := json.NewDecoder(r.Body)
    parameters := &params{}
    err := decoder.Decode(parameters) 
    if err != nil {
		log.Printf("Error marsheling json: %v", err)
    }

    new_chirp, err := cfg.db.CreateChirp(ReplaceProfaneWords(parameters.Body))
    if err != nil {
        log.Println(err)
    }

    RespondWithJSON(w, http.StatusCreated, Chirp{
        ID: new_chirp.ID,
        Body: new_chirp.Body,
    })
}
