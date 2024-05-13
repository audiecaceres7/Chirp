package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

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

    new_chirp, err := cfg.db.CreateChirp(parameters.Body)
    if err != nil {
        log.Println(err)
    }

    RespondWithJSON(w, http.StatusCreated, new_chirp)
}

func (cfg *apiConfig) HandleChirps(w http.ResponseWriter, r *http.Request) {
    all_chirps, err := cfg.db.GetChirps()
    if err != nil {
        log.Println(err)
    }

    RespondWithJSON(w, http.StatusCreated, all_chirps)
}

func RespondWithJSON(w http.ResponseWriter, code int, parameters interface{}) {
	data, err := json.Marshal(parameters)
	if err != nil {
		log.Printf("Error marsheling json: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	type error struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, &error{
		Error: message,
	})
}
