package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) HandleChirps(w http.ResponseWriter, r *http.Request) {
    all_chirps, err := cfg.db.GetChirps()
    if err != nil {
        log.Println(err)
    }

    RespondWithJSON(w, http.StatusOK, all_chirps)
}
