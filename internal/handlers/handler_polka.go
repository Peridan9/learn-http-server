package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type UserData struct {
		UserID string `json:"user_id"`
	}

	type Parameters struct {
		Event string   `json:"event"`
		Data  UserData `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := Parameters{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	if parameters.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Invalid event", nil)
		return
	}

	//convert userID to uuid
	userID, err := uuid.Parse(parameters.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid user ID", err)
		return
	}

	//update user in the database
	err = cfg.DB.UpgradeUserRedByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	//respond with 200 Ok and empty body
	respondWithJSON(w, http.StatusNoContent, nil)

}
