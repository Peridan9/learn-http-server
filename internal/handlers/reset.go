package handlers

import "net/http"

func (cfg *APIConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset only allowed in dev environment", nil)
		return
	}

	err := cfg.DB.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not reset users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset successful"))
}
