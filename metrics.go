package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(cfg.fileserverHits)
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    html := `
        <html>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>
    `
    w.Write([]byte(fmt.Sprintf(html, cfg.fileserverHits)))
}
