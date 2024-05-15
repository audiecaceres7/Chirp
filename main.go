package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"chirpy.com/database"
)

const (
	port          = "8080"
	root          = "."
	database_path = "./database/database.json"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	mux := http.NewServeMux()
    dbg := flag.Bool("debug", false, "Enable debug mode")
    flag.Parse()
    if *dbg {
        fmt.Println("debug mode enabled")
        err := os.Remove("database.json")
        if err != nil {
            fmt.Printf("Error deleting file: %v\n", err)
        }
    }

	db, err := database.NewDB(database_path)
	config := apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	err = config.db.EnsureDB()
	if err != nil {
		fmt.Println("Created database and ensured")
	}

	mux.HandleFunc("POST /api/login", config.HandlerLogin)
	mux.HandleFunc("POST /api/chirps", config.HandleChirp)
	mux.HandleFunc("POST /api/users", config.HandlerCreateUser)
	mux.HandleFunc("GET /api/chirps", config.HandleChirps)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", config.HandleChirpId)
	mux.HandleFunc("GET /api/healthz", HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", config.HandleMetrics)
	mux.HandleFunc("/api/reset", config.HandleReset)
	mux.Handle("/app/*", config.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	// starting server
	fmt.Printf("Serving files from %v on port: %v\n", root, port)
	log.Fatal(server.ListenAndServe())
}
