package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
    port = "8080"
    root = "."
)

type apiConfig struct {
	fileserverHits int
}

func main() {
    mux := http.NewServeMux()

    config := apiConfig{
        fileserverHits: 0,
    }

    mux.HandleFunc("GET /healthz", HandleReadiness)
    mux.HandleFunc("GET /metrics", config.HandleMetrics)
    mux.HandleFunc("/reset", config.HandleReset)
    mux.Handle("/app/*", config.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))

    server := &http.Server{
        Handler: mux,
        Addr: ":" + port,
    }

    // starting server
    fmt.Printf("Serving files from %v on port: %v\n", root, port)
    log.Fatal(server.ListenAndServe())
}
