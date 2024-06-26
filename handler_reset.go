package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    cfg.fileserverHits = 0
    w.Write([]byte(fmt.Sprintf("Hits have been reseted to %v", cfg.fileserverHits)))
}
