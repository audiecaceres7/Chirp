package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) HandleChirpId(w http.ResponseWriter, r *http.Request) {
	all_chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Printf("couldn't get chirps: %v", err)
	}

	id, err := strconv.Atoi(r.PathValue("chirp_id"));
	if err != nil {
		log.Printf("Coudln't conver ASCII to int: %v", err)
	}

    for _, chirp := range all_chirps {
        if chirp.ID == id {
            RespondWithJSON(w, http.StatusOK, Chirp{
                ID:   all_chirps[id-1].ID,
                Body: all_chirps[id-1].Body,
            })
            return
        }
    }

    RespondWithError(w, http.StatusNotFound, "Chirp does not exist")
}
