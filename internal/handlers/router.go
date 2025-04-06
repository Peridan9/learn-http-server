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
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdatePasswordAndEmail)

	// login
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	// Tokens
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	// chirps
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)

	return mux
}
