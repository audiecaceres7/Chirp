package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
