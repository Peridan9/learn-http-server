package handlers

import (
	"net/http"
)

func (cfg *APIConfig) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// app
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	// healthz
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// admin
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	// users
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)

	// login
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	// chirps
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)

	return mux
}
